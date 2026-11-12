package audit_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/audit"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/tenant"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/database/dbtest"
)

type retention int

func (r retention) Value(context.Context, pgx.Tx, string) (json.RawMessage, error) {
	return json.Marshal(int(r))
}

type allowAll struct{}

func (allowAll) AllowedInTenant(context.Context, pgx.Tx, auth.Principal, string) (bool, error) {
	return true, nil
}

func TestPrunerKeepsOnlyRetentionWindow(t *testing.T) {
	ctx := context.Background()
	pool := dbtest.New(t)
	if err := database.NewMigrator(pool).Apply(ctx, tenant.Migrations(), identity.Migrations(), audit.Migrations()); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	created, err := tenant.NewStore(pool).Create(ctx, "kaznu", "КазНУ")
	if err != nil {
		t.Fatalf("tenant: %v", err)
	}
	tx := database.NewTenantTx(pool)
	err = tx.Run(ctx, created.ID, func(tx pgx.Tx) error {
		if err := (audit.Log{}).Record(ctx, tx, "", audit.Change{Action: "core.setting.changed", TargetID: "fresh"}); err != nil {
			return err
		}
		return (audit.Log{}).Record(ctx, tx, "", audit.Change{Action: "core.setting.changed", TargetID: "stale"})
	})
	if err != nil {
		t.Fatalf("record: %v", err)
	}
	err = pgx.BeginFunc(ctx, pool, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, `SELECT set_config('lodestar.tenant', $1, true)`, created.ID); err != nil {
			return err
		}
		_, err := tx.Exec(ctx, `UPDATE core.audit_log SET at = now() - interval '40 days' WHERE target_id = 'stale'`)
		return err
	})
	if err != nil {
		t.Fatalf("age entry: %v", err)
	}

	if err := audit.NewPruner(pool, retention(30), time.Hour).PruneAll(ctx); err != nil {
		t.Fatalf("prune: %v", err)
	}

	entries, err := audit.NewService(tx, allowAll{}).List(ctx, auth.Principal{TenantID: created.ID}, audit.Filter{})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(entries) != 1 || entries[0].TargetID != "fresh" || entries[0].ActorID != nil {
		t.Fatalf("entries = %+v", entries)
	}
}
