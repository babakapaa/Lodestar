package study

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
)

func (s AssignmentService) SearchMap(ctx context.Context, p auth.Principal, id string) (SearchMapDraft, error) {
	var d SearchMapDraft
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		a, err := s.studied(ctx, tx, p, id, false)
		if err != nil {
			return err
		}
		m, at, err := s.maps.Get(ctx, tx, id, p.UserID)
		if err != nil {
			return err
		}
		d = draft(a, m, at)
		return nil
	})
	return d, err
}

func (s AssignmentService) SaveSearchMap(ctx context.Context, p auth.Principal, id string, in SearchMap) (SearchMapDraft, error) {
	m, err := in.normalized()
	if err != nil {
		return SearchMapDraft{}, err
	}
	var d SearchMapDraft
	err = s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		a, err := s.studied(ctx, tx, p, id, true)
		if err != nil {
			return err
		}
		at, err := s.maps.Put(ctx, tx, id, p.UserID, m)
		if err != nil {
			return err
		}
		d = draft(a, m, &at)
		return nil
	})
	return d, err
}

func (s AssignmentService) studied(ctx context.Context, tx pgx.Tx, p auth.Principal, id string, lock bool) (Assignment, error) {
	a, _, st, err := s.visible(ctx, tx, p, id, lock)
	if err != nil {
		return Assignment{}, err
	}
	if st.Role != RoleStudent {
		return Assignment{}, access.ErrForbidden
	}
	return a, nil
}

func draft(a Assignment, m SearchMap, at *time.Time) SearchMapDraft {
	rules := MapRules{}
	if a.SearchMap != nil {
		rules = *a.SearchMap
	}
	return SearchMapDraft{Map: m, Rules: a.SearchMap, Missing: m.Missing(rules), ChangedAt: at}
}
