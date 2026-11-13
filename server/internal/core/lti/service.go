package lti

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/audit"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/core/tenant"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

type Service struct {
	tenants   tenant.Store
	tx        database.TenantTx
	access    access.Service
	links     module.LinkSources
	keys      KeyStore
	platforms PlatformStore
	contexts  ContextStore
	members   MemberStore
	audit     audit.Log
	publicURL string
}

func (s Service) Tool(ctx context.Context, p auth.Principal) (ToolURLs, error) {
	t, err := s.tenants.ByID(ctx, p.TenantID)
	if err != nil {
		return ToolURLs{}, err
	}
	err = s.tx.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if err := s.require(ctx, tx, p); err != nil {
			return err
		}
		_, err := s.keys.Current(ctx, tx)
		return err
	})
	if err != nil {
		return ToolURLs{}, err
	}
	return NewToolURLs(s.publicURL, t.Slug), nil
}

func (s Service) Platforms(ctx context.Context, p auth.Principal) ([]Platform, error) {
	var list []Platform
	err := s.tx.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if err := s.require(ctx, tx, p); err != nil {
			return err
		}
		var err error
		list, err = s.platforms.List(ctx, tx)
		return err
	})
	return list, err
}

func (s Service) CreatePlatform(ctx context.Context, p auth.Principal, in PlatformInput) (Platform, error) {
	in, err := in.Normalized()
	if err != nil {
		return Platform{}, err
	}
	var out Platform
	err = s.tx.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if err := s.require(ctx, tx, p); err != nil {
			return err
		}
		var err error
		if out, err = s.platforms.Insert(ctx, tx, in); err != nil {
			return err
		}
		return s.audit.Record(ctx, tx, p.UserID, audit.Change{Action: ActionPlatformCreated, TargetKind: "lti.platform", TargetID: out.ID, Details: in})
	})
	return out, err
}

func (s Service) UpdatePlatform(ctx context.Context, p auth.Principal, id string, in PlatformInput) (Platform, error) {
	in, err := in.Normalized()
	if err != nil {
		return Platform{}, err
	}
	var out Platform
	err = s.tx.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if err := s.require(ctx, tx, p); err != nil {
			return err
		}
		var err error
		if out, err = s.platforms.Update(ctx, tx, id, in); err != nil {
			return err
		}
		return s.audit.Record(ctx, tx, p.UserID, audit.Change{Action: ActionPlatformUpdated, TargetKind: "lti.platform", TargetID: id, Details: in})
	})
	return out, err
}

func (s Service) DeletePlatform(ctx context.Context, p auth.Principal, id string) error {
	return s.tx.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if err := s.require(ctx, tx, p); err != nil {
			return err
		}
		before, err := s.platforms.ByID(ctx, tx, id)
		if err != nil {
			return err
		}
		if err := s.platforms.Delete(ctx, tx, id); err != nil {
			return err
		}
		return s.audit.Record(ctx, tx, p.UserID, audit.Change{Action: ActionPlatformDeleted, TargetKind: "lti.platform", TargetID: id, Details: before})
	})
}

func (s Service) Context(ctx context.Context, p auth.Principal, id string) (ContextView, error) {
	var v ContextView
	err := s.tx.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		c, m, err := s.member(ctx, tx, p, id)
		if err != nil {
			return err
		}
		platform, err := s.platforms.ByID(ctx, tx, c.PlatformID)
		if err != nil {
			return err
		}
		v = ContextView{ID: c.ID, Title: c.Title, Label: c.Label, Platform: platform.Name, Instructor: m.Instructor, Options: []module.LinkTarget{}}
		if v.Link, err = s.describe(ctx, tx, c); err != nil {
			return err
		}
		if !m.Instructor {
			return nil
		}
		for _, source := range s.links {
			options, err := source.Linkable(ctx, tx, p)
			if err != nil {
				return err
			}
			v.Options = append(v.Options, options...)
		}
		return nil
	})
	return v, err
}

func (s Service) Link(ctx context.Context, p auth.Principal, id string, in LinkInput) (ContextView, error) {
	source, ok := s.links.Find(in.Kind)
	if !ok {
		return ContextView{}, ErrUnknownLinkKind
	}
	err := s.tx.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		c, err := s.instructor(ctx, tx, p, id)
		if err != nil {
			return err
		}
		if err := source.CanLink(ctx, tx, p, in.ID); err != nil {
			return err
		}
		if err := s.contexts.Link(ctx, tx, c.ID, in.Kind, in.ID, p.UserID); err != nil {
			return err
		}
		return s.audit.Record(ctx, tx, p.UserID, audit.Change{Action: ActionContextLinked, TargetKind: "lti.context", TargetID: c.ID, Details: in})
	})
	if err != nil {
		return ContextView{}, err
	}
	return s.Context(ctx, p, id)
}

func (s Service) Unlink(ctx context.Context, p auth.Principal, id string) (ContextView, error) {
	err := s.tx.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		c, err := s.instructor(ctx, tx, p, id)
		if err != nil {
			return err
		}
		if c.OriginKind == nil {
			return nil
		}
		if err := s.contexts.Unlink(ctx, tx, c.ID); err != nil {
			return err
		}
		return s.audit.Record(ctx, tx, p.UserID, audit.Change{
			Action: ActionContextUnlinked, TargetKind: "lti.context", TargetID: c.ID,
			Details: LinkInput{Kind: *c.OriginKind, ID: *c.OriginID},
		})
	})
	if err != nil {
		return ContextView{}, err
	}
	return s.Context(ctx, p, id)
}

func (s Service) member(ctx context.Context, tx pgx.Tx, p auth.Principal, id string) (Context, Member, error) {
	c, err := s.contexts.ByID(ctx, tx, id)
	if err != nil {
		return Context{}, Member{}, err
	}
	m, err := s.members.Find(ctx, tx, c.ID, p.UserID)
	return c, m, err
}

func (s Service) instructor(ctx context.Context, tx pgx.Tx, p auth.Principal, id string) (Context, error) {
	c, m, err := s.member(ctx, tx, p, id)
	if err != nil {
		return Context{}, err
	}
	if !m.Instructor {
		return Context{}, access.ErrForbidden
	}
	return c, nil
}

func (s Service) describe(ctx context.Context, tx pgx.Tx, c Context) (*module.LinkTarget, error) {
	if c.OriginKind == nil {
		return nil, nil
	}
	source, ok := s.links.Find(*c.OriginKind)
	if !ok {
		return nil, nil
	}
	target, err := source.Describe(ctx, tx, *c.OriginID)
	if errors.Is(err, module.ErrLinkTargetNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &target, nil
}

func (s Service) require(ctx context.Context, tx pgx.Tx, p auth.Principal) error {
	return s.access.Require(ctx, tx, p, access.TenantScope(p.TenantID), ManageIntegrations)
}
