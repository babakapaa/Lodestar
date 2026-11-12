package database

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const migrationLockID = 7_310_455_201

const bootstrapSQL = `
CREATE TABLE IF NOT EXISTS lodestar_migrations (
	owner      text        NOT NULL,
	version    integer     NOT NULL,
	name       text        NOT NULL,
	applied_at timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (owner, version)
);
DO $$ BEGIN
	CREATE ROLE ` + AppRole + ` NOLOGIN;
EXCEPTION WHEN duplicate_object OR unique_violation THEN NULL;
END $$;`

type Migrator struct {
	pool *pgxpool.Pool
}

func NewMigrator(pool *pgxpool.Pool) Migrator {
	return Migrator{pool: pool}
}

func (m Migrator) Apply(ctx context.Context, sets ...Migrations) error {
	conn, err := m.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("acquire connection: %w", err)
	}
	defer conn.Release()

	if _, err := conn.Exec(ctx, "SELECT pg_advisory_lock($1)", migrationLockID); err != nil {
		return fmt.Errorf("lock migrations: %w", err)
	}
	defer func() {
		if _, err := conn.Exec(context.WithoutCancel(ctx), "SELECT pg_advisory_unlock($1)", migrationLockID); err != nil {
			slog.Error("unlock migrations", "error", err)
		}
	}()

	if _, err := conn.Exec(ctx, bootstrapSQL); err != nil {
		return fmt.Errorf("prepare migrations table: %w", err)
	}
	for _, set := range sets {
		if err := applySet(ctx, conn.Conn(), set); err != nil {
			return err
		}
	}
	return nil
}

func applySet(ctx context.Context, conn *pgx.Conn, set Migrations) error {
	files, err := readMigrationFiles(set.FS)
	if err != nil {
		return fmt.Errorf("%s: %w", set.Owner, err)
	}
	applied, err := appliedVersions(ctx, conn, set.Owner)
	if err != nil {
		return err
	}
	for _, file := range files {
		if _, done := applied[file.version]; done {
			continue
		}
		if err := applyFile(ctx, conn, set.Owner, file); err != nil {
			return err
		}
		slog.Info("migration applied", "owner", set.Owner, "name", file.name)
	}
	return nil
}

func appliedVersions(ctx context.Context, conn *pgx.Conn, owner string) (map[int]struct{}, error) {
	rows, err := conn.Query(ctx, "SELECT version FROM lodestar_migrations WHERE owner = $1", owner)
	if err != nil {
		return nil, fmt.Errorf("read applied migrations: %w", err)
	}
	versions, err := pgx.CollectRows(rows, pgx.RowTo[int])
	if err != nil {
		return nil, fmt.Errorf("read applied migrations: %w", err)
	}
	applied := make(map[int]struct{}, len(versions))
	for _, v := range versions {
		applied[v] = struct{}{}
	}
	return applied, nil
}

func applyFile(ctx context.Context, conn *pgx.Conn, owner string, file migrationFile) error {
	return pgx.BeginFunc(ctx, conn, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, file.sql); err != nil {
			return fmt.Errorf("%s/%s: %w", owner, file.name, err)
		}
		_, err := tx.Exec(ctx,
			"INSERT INTO lodestar_migrations (owner, version, name) VALUES ($1, $2, $3)",
			owner, file.version, file.name)
		if err != nil {
			return fmt.Errorf("record %s/%s: %w", owner, file.name, err)
		}
		return nil
	})
}
