package structure

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/access"
)

type Tree struct {
	units UnitStore
}

func (t Tree) Chain(ctx context.Context, tx pgx.Tx, tenantID string, s access.Scope) ([]access.Scope, error) {
	chain := []access.Scope{access.TenantScope(tenantID)}
	if s.IsTenant() {
		return chain, nil
	}
	if s.Kind != ScopeKind {
		return append(chain, s), nil
	}
	path, err := t.units.Path(ctx, tx, s.ID)
	if err != nil {
		return nil, err
	}
	if len(path) == 0 {
		return nil, fmt.Errorf("%w: %w", access.ErrInvalidScope, ErrUnitNotFound)
	}
	for _, u := range path {
		chain = append(chain, UnitScope(u.ID))
	}
	return chain, nil
}

func UnitScope(id string) access.Scope {
	return access.Scope{Kind: ScopeKind, ID: id}
}
