package audit

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

type Service struct {
	tenants database.TenantTx
	guard   Guard
	entries EntryStore
}

func NewService(tenants database.TenantTx, guard Guard) Service {
	return Service{tenants: tenants, guard: guard}
}

func (s Service) List(ctx context.Context, p auth.Principal, f Filter) ([]Entry, error) {
	var list []Entry
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		allowed, err := s.guard.AllowedInTenant(ctx, tx, p, ViewAudit)
		if err != nil {
			return err
		}
		if !allowed {
			return ErrForbidden
		}
		list, err = s.entries.Page(ctx, tx, f)
		return err
	})
	return list, err
}
