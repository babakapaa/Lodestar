package access

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type ScopeTree interface {
	Chain(ctx context.Context, tx pgx.Tx, tenantID string, s Scope) ([]Scope, error)
}

type FlatTree struct{}

func (FlatTree) Chain(_ context.Context, _ pgx.Tx, tenantID string, s Scope) ([]Scope, error) {
	if s.IsTenant() {
		return []Scope{TenantScope(tenantID)}, nil
	}
	return []Scope{TenantScope(tenantID), s}, nil
}
