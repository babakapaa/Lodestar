package settings

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type ModuleStore struct{}

func (ModuleStore) Disabled(ctx context.Context, tx pgx.Tx) (map[string]bool, error) {
	rows, err := tx.Query(ctx, `SELECT module FROM core.tenant_modules WHERE NOT enabled`)
	if err != nil {
		return nil, fmt.Errorf("read modules: %w", err)
	}
	names, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return nil, fmt.Errorf("read modules: %w", err)
	}
	off := make(map[string]bool, len(names))
	for _, n := range names {
		off[n] = true
	}
	return off, nil
}

func (ModuleStore) Enabled(ctx context.Context, tx pgx.Tx, name string) (bool, error) {
	var enabled bool
	err := tx.QueryRow(ctx, `SELECT enabled FROM core.tenant_modules WHERE module = $1`, name).Scan(&enabled)
	if errors.Is(err, pgx.ErrNoRows) {
		return true, nil
	}
	return enabled, err
}

func (ModuleStore) Set(ctx context.Context, tx pgx.Tx, name string, enabled bool, by string) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO core.tenant_modules (module, enabled, updated_by) VALUES ($1, $2, $3)
		ON CONFLICT (tenant_id, module) DO UPDATE
		SET enabled = EXCLUDED.enabled, updated_by = EXCLUDED.updated_by, updated_at = now()`,
		name, enabled, by)
	if err != nil {
		return fmt.Errorf("save module state: %w", err)
	}
	return nil
}
