package chat

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

const syncEvery = time.Minute

type RosterWorker struct {
	pool     *pgxpool.Pool
	tenants  database.TenantTx
	settings module.Settings
	rosters  module.Rosters
	sync     RosterSync
	every    time.Duration
}

func NewRosterWorker(k module.Kernel) RosterWorker {
	return RosterWorker{
		pool: k.Pool, tenants: k.Tenants, settings: k.Settings, rosters: k.Rosters,
		sync: newRosterSync(k), every: syncEvery,
	}
}

func (x RosterWorker) Run(ctx context.Context) {
	ticker := time.NewTicker(x.every)
	defer ticker.Stop()
	for {
		if err := x.SyncAll(ctx); err != nil && ctx.Err() == nil {
			slog.Warn("sync course chats", "error", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (x RosterWorker) SyncAll(ctx context.Context) error {
	rows, err := x.pool.Query(ctx, `SELECT id::text FROM core.tenants`)
	if err != nil {
		return err
	}
	tenants, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return err
	}
	for _, tenantID := range tenants {
		if err := x.SyncTenant(ctx, tenantID); err != nil {
			return err
		}
	}
	return nil
}

func (x RosterWorker) SyncTenant(ctx context.Context, tenantID string) error {
	enabled, err := x.settings.ModuleEnabled(ctx, tenantID, Name)
	if err != nil || !enabled {
		return err
	}
	for _, source := range x.rosters {
		owner, _, _ := strings.Cut(source.Kind(), ".")
		on, err := x.settings.ModuleEnabled(ctx, tenantID, owner)
		if err != nil {
			return err
		}
		if !on {
			continue
		}
		var all []module.Roster
		err = x.tenants.Run(ctx, tenantID, func(tx pgx.Tx) error {
			var err error
			all, err = source.All(ctx, tx)
			return err
		})
		if err != nil {
			return err
		}
		for _, r := range all {
			err := x.tenants.Run(ctx, tenantID, func(tx pgx.Tx) error {
				return x.sync.Sync(ctx, tx, tenantID, r)
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
