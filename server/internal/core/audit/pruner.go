package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

type Pruner struct {
	pool     *pgxpool.Pool
	tenants  database.TenantTx
	settings Settings
	every    time.Duration
	entries  EntryStore
}

func NewPruner(pool *pgxpool.Pool, settings Settings, every time.Duration) Pruner {
	return Pruner{pool: pool, tenants: database.NewTenantTx(pool), settings: settings, every: every}
}

func (p Pruner) Run(ctx context.Context) {
	ticker := time.NewTicker(p.every)
	defer ticker.Stop()
	for {
		if err := p.PruneAll(ctx); err != nil && ctx.Err() == nil {
			slog.Warn("prune audit log", "error", err)
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
			raw, err := p.settings.Value(ctx, tx, RetentionSetting)
			if err != nil {
				return err
			}
			var days int
			if err := json.Unmarshal(raw, &days); err != nil {
				return fmt.Errorf("read %s: %w", RetentionSetting, err)
			}
			return p.entries.Prune(ctx, tx, time.Duration(days)*24*time.Hour)
		})
		if err != nil {
			return err
		}
	}
	return nil
}
