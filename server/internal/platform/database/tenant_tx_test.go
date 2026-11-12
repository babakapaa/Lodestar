package database_test

import (
	"context"
	"testing"
	"testing/fstest"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/database/dbtest"
)

const isolatedTableSQL = `
CREATE TABLE notes (
	tenant_id uuid NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid,
	body      text NOT NULL
);
ALTER TABLE notes ENABLE ROW LEVEL SECURITY;
ALTER TABLE notes FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON notes
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);
GRANT SELECT, INSERT ON notes TO lodestar_app;`

const (
	tenantA = "11111111-1111-1111-1111-111111111111"
	tenantB = "22222222-2222-2222-2222-222222222222"
)

func TestTenantTxIsolatesTenants(t *testing.T) {
	pool := dbtest.New(t)
	ctx := context.Background()
	set := database.Migrations{Owner: "isolation", FS: fstest.MapFS{
		"0001_notes.sql": {Data: []byte(isolatedTableSQL)},
	}}
	if err := database.NewMigrator(pool).Apply(ctx, set); err != nil {
		t.Fatalf("Apply: %v", err)
	}
	tenants := database.NewTenantTx(pool)

	insert := func(tenant, body string) {
		err := tenants.Run(ctx, tenant, func(tx pgx.Tx) error {
			_, err := tx.Exec(ctx, "INSERT INTO notes (body) VALUES ($1)", body)
			return err
		})
		if err != nil {
			t.Fatalf("insert for %s: %v", tenant, err)
		}
	}
	insert(tenantA, "a1")
	insert(tenantA, "a2")
	insert(tenantB, "b1")

	var seen int
	err := tenants.Run(ctx, tenantB, func(tx pgx.Tx) error {
		return tx.QueryRow(ctx, "SELECT count(*) FROM notes").Scan(&seen)
	})
	if err != nil {
		t.Fatalf("count: %v", err)
	}
	if seen != 1 {
		t.Fatalf("tenant B sees %d rows, want 1", seen)
	}
}

func TestTenantTxBlocksWritesIntoAnotherTenant(t *testing.T) {
	pool := dbtest.New(t)
	ctx := context.Background()
	set := database.Migrations{Owner: "isolation", FS: fstest.MapFS{
		"0001_notes.sql": {Data: []byte(isolatedTableSQL)},
	}}
	if err := database.NewMigrator(pool).Apply(ctx, set); err != nil {
		t.Fatalf("Apply: %v", err)
	}

	err := database.NewTenantTx(pool).Run(ctx, tenantA, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, "INSERT INTO notes (tenant_id, body) VALUES ($1, 'x')", tenantB)
		return err
	})
	if err == nil {
		t.Fatal("tenant A wrote a row owned by tenant B")
	}
}
