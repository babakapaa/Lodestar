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
	deepLinks DeepLinkStore
	values    Values
	sync      SyncService
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
		c, err := s.contexts.ByID(ctx, tx, id)
		if err != nil {
			return err
		}
		v, err = s.view(ctx, tx, p, c)
		return err
	})
	return v, err
}

func (s Service) ByOrigin(ctx context.Context, p auth.Principal, kind, originID string) ([]ContextView, error) {
	source, ok := s.links.Find(kind)
	if !ok {
		return nil, ErrUnknownLinkKind
	}
	list := []ContextView{}
	err := s.tx.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if err := source.CanLink(ctx, tx, p, originID); err != nil {
			if !errors.Is(err, access.ErrForbidden) {
				return err
			}
			if err := source.CanManage(ctx, tx, p, originID); err != nil {
				return err
			}
		}
		contexts, err := s.contexts.ByOrigin(ctx, tx, kind, originID)
		if err != nil {
			return err
		}
		for _, c := range contexts {
			v, err := s.view(ctx, tx, p, c)
			if err != nil {
				return err
			}
			list = append(list, v)
		}
		return nil
	})
	return list, err
}

func (s Service) SetSource(ctx context.Context, p auth.Principal, id string, in SourceInput) (ContextView, error) {
	if !in.valid() {
		return ContextView{}, ErrInvalidSource
	}
	err := s.tx.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		c, err := s.contexts.ByID(ctx, tx, id)
		if err != nil {
			return err
		}
		tenant, err := tenantSource(ctx, tx, s.values)
		if err != nil {
			return err
		}
		if tenant != SourceUnits {
			return ErrSourceFixed
		}
		manages, err := s.manages(ctx, tx, p, c)
		if err != nil {
			return err
		}
		if !manages {
			return access.ErrForbidden
		}
		if err := s.contexts.SetSource(ctx, tx, c.ID, in.Source); err != nil {
			return err
		}
		return s.audit.Record(ctx, tx, p.UserID, audit.Change{Action: ActionSourceChanged, TargetKind: "lti.context", TargetID: c.ID, Details: in})
	})
	if err != nil {
		return ContextView{}, err
	}
	return s.Context(ctx, p, id)
}

func (s Service) Sync(ctx context.Context, p auth.Principal, id string) (SyncReport, error) {
	err := s.tx.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		c, err := s.contexts.ByID(ctx, tx, id)
		if err != nil {
			return err
		}
		v, err := s.view(ctx, tx, p, c)
		if err != nil {
			return err
		}
		if !v.CanSync {
			return access.ErrForbidden
		}
		return nil
	})
	if err != nil {
		return SyncReport{}, err
	}
	return s.sync.Sync(ctx, p.TenantID, id)
}

func (s Service) view(ctx context.Context, tx pgx.Tx, p auth.Principal, c Context) (ContextView, error) {
	m, err := s.members.Find(ctx, tx, c.ID, p.UserID)
	member := err == nil
	if err != nil && !errors.Is(err, ErrContextNotFound) {
		return ContextView{}, err
	}
	manages, err := s.manages(ctx, tx, p, c)
	if err != nil {
		return ContextView{}, err
	}
	if !member && !manages {
		return ContextView{}, ErrContextNotFound
	}
	platform, err := s.platforms.ByID(ctx, tx, c.PlatformID)
	if err != nil {
		return ContextView{}, err
	}
	tenant, err := tenantSource(ctx, tx, s.values)
	if err != nil {
		return ContextView{}, err
	}
	v := ContextView{
		ID: c.ID, Title: c.Title, Label: c.Label, Platform: platform.Name, Instructor: m.Instructor,
		Options: []module.LinkTarget{}, Source: effectiveSource(tenant, c),
		SyncedAt: c.SyncedAt, SyncError: c.SyncError,
	}
	if v.Link, err = s.describe(ctx, tx, c); err != nil {
		return ContextView{}, err
	}
	linked := v.Link != nil
	v.CanChoose = linked && manages && tenant == SourceUnits
	v.CanSync = linked && (m.Instructor || manages)
	if !m.Instructor {
		return v, nil
	}
	for _, source := range s.links {
		options, err := source.Linkable(ctx, tx, p)
		if err != nil {
			return ContextView{}, err
		}
		v.Options = append(v.Options, options...)
	}
	return v, nil
}

func (s Service) manages(ctx context.Context, tx pgx.Tx, p auth.Principal, c Context) (bool, error) {
	if c.OriginKind == nil {
		return false, nil
	}
	source, ok := s.links.Find(*c.OriginKind)
	if !ok {
		return false, nil
	}
	err := source.CanManage(ctx, tx, p, *c.OriginID)
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, access.ErrForbidden), errors.Is(err, module.ErrLinkTargetNotFound):
		return false, nil
	default:
		return false, err
	}
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

func (s Service) instructor(ctx context.Context, tx pgx.Tx, p auth.Principal, id string) (Context, error) {
	c, err := s.contexts.ByID(ctx, tx, id)
	if err != nil {
		return Context{}, err
	}
	m, err := s.members.Find(ctx, tx, c.ID, p.UserID)
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
