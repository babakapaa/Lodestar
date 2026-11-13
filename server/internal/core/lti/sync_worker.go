package lti

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SyncWorker struct {
	pool     *pgxpool.Pool
	sync     SyncService
	contexts ContextStore
	every    time.Duration
}

func (w SyncWorker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.every)
	defer ticker.Stop()
	for {
		if err := w.SyncAll(ctx); err != nil && ctx.Err() == nil {
			slog.Warn("lti sync", "error", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (w SyncWorker) SyncAll(ctx context.Context) error {
	rows, err := w.pool.Query(ctx, `SELECT id::text FROM core.tenants`)
	if err != nil {
		return err
	}
	tenants, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return err
	}
	for _, tenantID := range tenants {
		var linked []Context
		err := w.sync.tx.Run(ctx, tenantID, func(tx pgx.Tx) error {
			var err error
			linked, err = w.contexts.Linked(ctx, tx)
			return err
		})
		if err != nil {
			return err
		}
		for _, c := range linked {
			if _, err := w.sync.Sync(ctx, tenantID, c.ID); err != nil {
				if ctx.Err() != nil {
					return ctx.Err()
				}
				slog.Warn("lti context sync failed", "tenant", tenantID, "context", c.ID, "error", err)
			}
		}
	}
	return nil
}
