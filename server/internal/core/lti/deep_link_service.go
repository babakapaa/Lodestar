package lti

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/module"
)

const maxDeepLinkItems = 20

func (s Service) DeepLink(ctx context.Context, p auth.Principal, id string) (DeepLink, error) {
	var d DeepLink
	err := s.tx.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		var err error
		if d, err = s.deepLinks.Open(ctx, tx, id, p.UserID); err != nil {
			return err
		}
		d.Targets = []module.LinkTarget{}
		for _, source := range s.links {
			list, err := source.Linkable(ctx, tx, p)
			if err != nil {
				return err
			}
			d.Targets = append(d.Targets, list...)
		}
		return nil
	})
	return d, err
}

func (s Service) AnswerDeepLink(ctx context.Context, p auth.Principal, id string, in DeepLinkInput) (DeepLinkAnswer, error) {
	if len(in.Items) == 0 || len(in.Items) > maxDeepLinkItems {
		return DeepLinkAnswer{}, ErrInvalidDeepLink
	}
	t, err := s.tenants.ByID(ctx, p.TenantID)
	if err != nil {
		return DeepLinkAnswer{}, err
	}
	code, err := randomString()
	if err != nil {
		return DeepLinkAnswer{}, err
	}
	err = s.tx.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		d, err := s.deepLinks.Open(ctx, tx, id, p.UserID)
		if err != nil {
			return err
		}
		if !d.AcceptMultiple && len(in.Items) > 1 {
			return ErrInvalidDeepLink
		}
		platform, err := s.platforms.ByID(ctx, tx, d.PlatformID)
		if err != nil {
			return err
		}
		targets := make([]module.LinkTarget, 0, len(in.Items))
		for _, item := range in.Items {
			source, ok := s.links.Find(item.Kind)
			if !ok {
				return ErrUnknownLinkKind
			}
			if err := source.CanLink(ctx, tx, p, item.ID); err != nil {
				return err
			}
			target, err := source.Describe(ctx, tx, item.ID)
			if err != nil {
				return err
			}
			targets = append(targets, target)
		}
		key, err := s.keys.Current(ctx, tx)
		if err != nil {
			return err
		}
		response, err := deepLinkResponse(platform, d, targets, key)
		if err != nil {
			return err
		}
		return s.deepLinks.Answer(ctx, tx, d.ID, response, code)
	})
	if err != nil {
		return DeepLinkAnswer{}, err
	}
	return DeepLinkAnswer{FormURL: deepLinkReturnURL(s.publicURL, t.Slug, code)}, nil
}
