package files

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

type Sweeper struct {
	pool    *pgxpool.Pool
	tenants database.TenantTx
	storage Storage
	after   time.Duration
	every   time.Duration
	store   FileStore
}

func NewSweeper(pool *pgxpool.Pool, storage Storage, after, every time.Duration) Sweeper {
	return Sweeper{pool: pool, tenants: database.NewTenantTx(pool), storage: storage, after: after, every: every}
}

func (s Sweeper) Run(ctx context.Context) {
	ticker := time.NewTicker(s.every)
	defer ticker.Stop()
	for {
		if err := s.SweepAll(ctx); err != nil && ctx.Err() == nil {
			slog.Warn("sweep abandoned uploads", "error", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (s Sweeper) SweepAll(ctx context.Context) error {
	rows, err := s.pool.Query(ctx, `SELECT id::text FROM core.tenants`)
	if err != nil {
		return err
	}
	ids, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return err
	}
	for _, tenantID := range ids {
		err := s.tenants.Run(ctx, tenantID, func(tx pgx.Tx) error {
			abandoned, err := s.store.Abandoned(ctx, tx, s.after)
			if err != nil {
				return err
			}
			for _, fileID := range abandoned {
				if err := s.storage.Delete(ctx, objectKey(tenantID, fileID)); err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}
