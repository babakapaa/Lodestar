package database_test

import (
	"context"
	"testing"
	"testing/fstest"

	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/database/dbtest"
)

func TestMigratorAppliesOnceInOrder(t *testing.T) {
	pool := dbtest.New(t)
	ctx := context.Background()
	set := database.Migrations{
		Owner: "sample",
		FS: fstest.MapFS{
			"0002_add_title.sql": {Data: []byte("ALTER TABLE sample_items ADD COLUMN title text")},
			"0001_create.sql":    {Data: []byte("CREATE TABLE sample_items (id int PRIMARY KEY)")},
			"README.md":          {Data: []byte("ignored")},
		},
	}
	migrator := database.NewMigrator(pool)

	for range 2 {
		if err := migrator.Apply(ctx, set); err != nil {
			t.Fatalf("Apply: %v", err)
		}
	}

	var count int
	if err := pool.QueryRow(ctx, "SELECT count(*) FROM lodestar_migrations WHERE owner = 'sample'").Scan(&count); err != nil {
		t.Fatalf("count: %v", err)
	}
	if count != 2 {
		t.Fatalf("applied %d migrations, want 2", count)
	}
	if _, err := pool.Exec(ctx, "INSERT INTO sample_items (id, title) VALUES (1, 'x')"); err != nil {
		t.Fatalf("schema not migrated: %v", err)
	}
}

func TestMigratorRollsBackFailedMigration(t *testing.T) {
	pool := dbtest.New(t)
	ctx := context.Background()
	set := database.Migrations{
		Owner: "broken",
		FS: fstest.MapFS{
			"0001_bad.sql": {Data: []byte("CREATE TABLE half_done (id int); SELECT missing_column FROM half_done")},
		},
	}

	if err := database.NewMigrator(pool).Apply(ctx, set); err == nil {
		t.Fatal("Apply succeeded on a broken migration")
	}

	var exists bool
	if err := pool.QueryRow(ctx, "SELECT to_regclass('half_done') IS NOT NULL").Scan(&exists); err != nil {
		t.Fatalf("check table: %v", err)
	}
	if exists {
		t.Fatal("broken migration left a table behind")
	}
}

func TestMigratorRejectsDuplicateVersions(t *testing.T) {
	pool := dbtest.New(t)
	set := database.Migrations{
		Owner: "dup",
		FS: fstest.MapFS{
			"0001_a.sql": {Data: []byte("SELECT 1")},
			"0001_b.sql": {Data: []byte("SELECT 1")},
		},
	}

	if err := database.NewMigrator(pool).Apply(context.Background(), set); err == nil {
		t.Fatal("Apply accepted two migrations with the same version")
	}
}
