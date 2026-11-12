package live

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

type Pruner struct {
	pool    *pgxpool.Pool
	tenants database.TenantTx
	keep    time.Duration
	every   time.Duration
	store   Store
}

func NewPruner(pool *pgxpool.Pool, keep, every time.Duration) Pruner {
	return Pruner{pool: pool, tenants: database.NewTenantTx(pool), keep: keep, every: every}
}

func (p Pruner) Run(ctx context.Context) {
	ticker := time.NewTicker(p.every)
	defer ticker.Stop()
	for {
		if err := p.PruneAll(ctx); err != nil && ctx.Err() == nil {
			slog.Warn("prune live events", "error", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (p Pruner) PruneAll(ctx context.Context) error {
	rows, err := p.pool.Query(ctx, `SELECT id::text FROM core.tenants`)
	if err != nil {
		return err
	}
	ids, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return err
	}
	for _, id := range ids {
		err := p.tenants.Run(ctx, id, func(tx pgx.Tx) error {
			return p.store.Prune(ctx, tx, p.keep)
		})
		if err != nil {
			return err
		}
	}
	return nil
}
