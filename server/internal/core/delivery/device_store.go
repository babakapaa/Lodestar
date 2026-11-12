package delivery

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

const deviceColumns = `id::text, platform, name, language, created_at, last_seen_at`

type DeviceStore struct{}

func (DeviceStore) Upsert(ctx context.Context, tx pgx.Tx, userID string, in DeviceInput) (Device, error) {
	row := tx.QueryRow(ctx, `
		INSERT INTO core.push_devices (user_id, platform, token, name, language) VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (tenant_id, token) DO UPDATE
		SET user_id = EXCLUDED.user_id, platform = EXCLUDED.platform, name = EXCLUDED.name,
		    language = EXCLUDED.language, last_seen_at = now()
		RETURNING `+deviceColumns, userID, in.Platform, in.Token, in.Name, in.Language)
	d, err := scanDevice(row)
	if err != nil {
		return Device{}, fmt.Errorf("save device: %w", err)
	}
	return d, nil
}

func (DeviceStore) List(ctx context.Context, tx pgx.Tx, userID string) ([]Device, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+deviceColumns+` FROM core.push_devices WHERE user_id::text = $1 ORDER BY last_seen_at DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("list devices: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (Device, error) {
		return scanDevice(row)
	})
}

func (DeviceStore) Targets(ctx context.Context, tx pgx.Tx, userID string) ([]target, error) {
	rows, err := tx.Query(ctx, `
		SELECT id::text, token, language FROM core.push_devices WHERE user_id::text = $1 ORDER BY last_seen_at DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("read devices: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (target, error) {
		var t target
		err := row.Scan(&t.id, &t.token, &t.language)
		return t, err
	})
}

func (DeviceStore) Delete(ctx context.Context, tx pgx.Tx, userID, id string) (bool, error) {
	tag, err := tx.Exec(ctx, `DELETE FROM core.push_devices WHERE id::text = $1 AND user_id::text = $2`, id, userID)
	return tag.RowsAffected() == 1, err
}

func (DeviceStore) Forget(ctx context.Context, tx pgx.Tx, id string) error {
	_, err := tx.Exec(ctx, `DELETE FROM core.push_devices WHERE id::text = $1`, id)
	return err
}

func scanDevice(row pgx.Row) (Device, error) {
	var d Device
	err := row.Scan(&d.ID, &d.Platform, &d.Name, &d.Language, &d.CreatedAt, &d.LastSeenAt)
	return d, err
}
