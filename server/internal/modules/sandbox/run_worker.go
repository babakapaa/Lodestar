package sandbox

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

const scheduleEvery = 15 * time.Second

type RunWorker struct {
	pool     *pgxpool.Pool
	tenants  database.TenantTx
	settings module.Settings
	service  Service
	every    time.Duration
}

func NewRunWorker(k module.Kernel) RunWorker {
	return RunWorker{pool: k.Pool, tenants: k.Tenants, settings: k.Settings, service: NewService(k), every: scheduleEvery}
}

func (w RunWorker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.every)
	defer ticker.Stop()
	for {
		if err := w.ScheduleAll(ctx); err != nil && ctx.Err() == nil {
			slog.Warn("schedule sandbox runs", "error", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (w RunWorker) ScheduleAll(ctx context.Context) error {
	rows, err := w.pool.Query(ctx, `SELECT id::text FROM core.tenants`)
	if err != nil {
		return err
	}
	tenants, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return err
	}
	for _, tenantID := range tenants {
		enabled, err := w.settings.ModuleEnabled(ctx, tenantID, Name)
		if err != nil {
			return err
		}
		if !enabled {
			continue
		}
		if err := w.tenants.Run(ctx, tenantID, func(tx pgx.Tx) error { return w.service.Schedule(ctx, tx) }); err != nil {
			return err
		}
	}
	return nil
}
