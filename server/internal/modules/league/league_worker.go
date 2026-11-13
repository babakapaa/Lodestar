package league

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

const syncEvery = time.Minute

type LeagueWorker struct {
	pool     *pgxpool.Pool
	tenants  database.TenantTx
	settings module.Settings
	service  Service
	every    time.Duration
}

func NewLeagueWorker(k module.Kernel) LeagueWorker {
	return LeagueWorker{pool: k.Pool, tenants: k.Tenants, settings: k.Settings, service: NewService(k), every: syncEvery}
}

func (w LeagueWorker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.every)
	defer ticker.Stop()
	for {
		if err := w.SyncAll(ctx); err != nil && ctx.Err() == nil {
			slog.Warn("sync leagues", "error", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (w LeagueWorker) SyncAll(ctx context.Context) error {
	rows, err := w.pool.Query(ctx, `SELECT id::text FROM core.tenants`)
	if err != nil {
		return err
	}
	tenants, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return err
	}
	for _, tenantID := range tenants {
		if err := w.SyncTenant(ctx, tenantID); err != nil {
			return err
		}
	}
	return nil
}

func (w LeagueWorker) SyncTenant(ctx context.Context, tenantID string) error {
	enabled, err := w.settings.ModuleEnabled(ctx, tenantID, Name)
	if err != nil || !enabled {
		return err
	}
	for _, source := range w.service.scores {
		var all []module.Roster
		err := w.tenants.Run(ctx, tenantID, func(tx pgx.Tx) error {
			var err error
			all, err = source.All(ctx, tx)
			return err
		})
		if err != nil {
			return err
		}
		for _, r := range all {
			if r.Archived {
				continue
			}
			err := w.tenants.Run(ctx, tenantID, func(tx pgx.Tx) error {
				_, err := w.service.sync(ctx, tx, tenantID, r, source)
				return err
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
