package profile

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type PreferenceStore struct{}

func (PreferenceStore) All(ctx context.Context, tx pgx.Tx, userID string) (map[string]json.RawMessage, error) {
	rows, err := tx.Query(ctx, `SELECT key, value FROM core.preferences WHERE user_id::text = $1`, userID)
	if err != nil {
		return nil, fmt.Errorf("read preferences: %w", err)
	}
	values := map[string]json.RawMessage{}
	var key string
	var value []byte
	_, err = pgx.ForEachRow(rows, []any{&key, &value}, func() error {
		values[key] = append(json.RawMessage(nil), value...)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("read preferences: %w", err)
	}
	return values, nil
}

func (PreferenceStore) Put(ctx context.Context, tx pgx.Tx, userID, key string, value json.RawMessage) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO core.preferences (user_id, key, value) VALUES ($1, $2, $3)
		ON CONFLICT (user_id, key) DO UPDATE SET value = EXCLUDED.value, updated_at = now()`,
		userID, key, []byte(value))
	if err != nil {
		return fmt.Errorf("save preference: %w", err)
	}
	return nil
}

func (PreferenceStore) Delete(ctx context.Context, tx pgx.Tx, userID, key string) error {
	_, err := tx.Exec(ctx, `DELETE FROM core.preferences WHERE user_id::text = $1 AND key = $2`, userID, key)
	return err
}

func (PreferenceStore) Get(ctx context.Context, tx pgx.Tx, userID, key string) (json.RawMessage, bool, error) {
	var value []byte
	err := tx.QueryRow(ctx, `SELECT value FROM core.preferences WHERE user_id::text = $1 AND key = $2`, userID, key).Scan(&value)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("read preference: %w", err)
	}
	return value, true, nil
}
