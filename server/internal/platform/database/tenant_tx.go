package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const AppRole = "lodestar_app"

const TenantSetting = "lodestar.tenant"

type TenantTx struct {
	pool *pgxpool.Pool
}

func NewTenantTx(pool *pgxpool.Pool) TenantTx {
	return TenantTx{pool: pool}
}

func (t TenantTx) Run(ctx context.Context, tenantID string, fn func(pgx.Tx) error) error {
	return pgx.BeginFunc(ctx, t.pool, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, "SET LOCAL ROLE "+AppRole); err != nil {
			return fmt.Errorf("switch to app role: %w", err)
		}
		if _, err := tx.Exec(ctx, "SELECT set_config($1, $2, true)", TenantSetting, tenantID); err != nil {
			return fmt.Errorf("set tenant: %w", err)
		}
		return fn(tx)
	})
}
