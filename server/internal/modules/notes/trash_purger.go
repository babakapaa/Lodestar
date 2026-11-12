package notes

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

const purgeEvery = time.Hour

type TrashPurger struct {
	pool     *pgxpool.Pool
	tenants  database.TenantTx
	settings module.Settings
	notes    NoteStore
	every    time.Duration
	now      func() time.Time
}

func NewTrashPurger(k module.Kernel) TrashPurger {
	return TrashPurger{pool: k.Pool, tenants: k.Tenants, settings: k.Settings, every: purgeEvery, now: time.Now}
}

func (t TrashPurger) Run(ctx context.Context) {
	ticker := time.NewTicker(t.every)
	defer ticker.Stop()
	for {
		if err := t.PurgeAll(ctx); err != nil && ctx.Err() == nil {
			slog.Warn("purge notes trash", "error", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (t TrashPurger) PurgeAll(ctx context.Context) error {
	rows, err := t.pool.Query(ctx, `SELECT id::text FROM core.tenants`)
	if err != nil {
		return err
	}
	tenants, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return err
	}
	for _, tenantID := range tenants {
		enabled, err := t.settings.ModuleEnabled(ctx, tenantID, Name)
		if err != nil || !enabled {
			continue
		}
		if err := t.tenants.Run(ctx, tenantID, func(tx pgx.Tx) error {
			raw, err := t.settings.Value(ctx, tx, TrashDaysSetting)
			if err != nil {
				return err
			}
			var days int
			if err := json.Unmarshal(raw, &days); err != nil {
				return fmt.Errorf("decode trash days: %w", err)
			}
			_, err = t.notes.Purge(ctx, tx, t.now().AddDate(0, 0, -days))
			return err
		}); err != nil {
			return err
		}
	}
	return nil
}
