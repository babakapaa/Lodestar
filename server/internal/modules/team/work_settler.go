package team

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/babakapaa/lodestar/server/internal/core/module"
)

const settleEvery = time.Minute

type WorkSettler struct {
	pool     *pgxpool.Pool
	settings module.Settings
	service  Service
	every    time.Duration
	now      func() time.Time
}

func NewWorkSettler(k module.Kernel) WorkSettler {
	return WorkSettler{pool: k.Pool, settings: k.Settings, service: NewService(k), every: settleEvery, now: time.Now}
}

func (w WorkSettler) Run(ctx context.Context) {
	ticker := time.NewTicker(w.every)
	defer ticker.Stop()
	for {
		if err := w.SettleAll(ctx); err != nil && ctx.Err() == nil {
			slog.Warn("settle team work", "error", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (w WorkSettler) SettleAll(ctx context.Context) error {
	return w.SettleAt(ctx, w.now())
}

func (w WorkSettler) SettleAt(ctx context.Context, now time.Time) error {
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
		if err := w.service.SettleDue(ctx, tenantID, now); err != nil {
			return err
		}
	}
	return nil
}
