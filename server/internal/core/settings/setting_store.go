package settings

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type SettingStore struct{}

func (SettingStore) All(ctx context.Context, tx pgx.Tx) (map[string]json.RawMessage, error) {
	rows, err := tx.Query(ctx, `SELECT key, value FROM core.tenant_settings`)
	if err != nil {
		return nil, fmt.Errorf("read settings: %w", err)
	}
	values := map[string]json.RawMessage{}
	var key string
	var value []byte
	_, err = pgx.ForEachRow(rows, []any{&key, &value}, func() error {
		values[key] = append(json.RawMessage(nil), value...)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("read settings: %w", err)
	}
	return values, nil
}

func (SettingStore) Put(ctx context.Context, tx pgx.Tx, key string, value json.RawMessage, by string) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO core.tenant_settings (key, value, updated_by) VALUES ($1, $2, $3)
		ON CONFLICT (tenant_id, key) DO UPDATE
		SET value = EXCLUDED.value, updated_by = EXCLUDED.updated_by, updated_at = now()`,
		key, []byte(value), by)
	if err != nil {
		return fmt.Errorf("save setting: %w", err)
	}
	return nil
}

func (SettingStore) Delete(ctx context.Context, tx pgx.Tx, key string) error {
	_, err := tx.Exec(ctx, `DELETE FROM core.tenant_settings WHERE key = $1`, key)
	return err
}

func (SettingStore) Get(ctx context.Context, tx pgx.Tx, key string) (json.RawMessage, bool, error) {
	var value []byte
	err := tx.QueryRow(ctx, `SELECT value FROM core.tenant_settings WHERE key = $1`, key).Scan(&value)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("read setting: %w", err)
	}
	return value, true, nil
}
