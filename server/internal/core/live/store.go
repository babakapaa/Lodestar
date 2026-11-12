package live

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

const replayLimit = 500

type Store struct{}

func (Store) Publish(ctx context.Context, tx pgx.Tx, tenantID string, userIDs []string, eventType string, payload any) error {
	if len(userIDs) == 0 {
		return nil
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("encode live payload: %w", err)
	}
	rows, err := tx.Query(ctx, `
		INSERT INTO core.live_events (user_id, type, payload)
		SELECT u, $2, $3 FROM unnest($1::uuid[]) AS u
		RETURNING id, user_id::text`, userIDs, eventType, data)
	if err != nil {
		return fmt.Errorf("publish live event: %w", err)
	}
	var signals []string
	var id int64
	var userID string
	_, err = pgx.ForEachRow(rows, []any{&id, &userID}, func() error {
		signals = append(signals, signal{tenantID: tenantID, userID: userID, eventID: id}.encode())
		return nil
	})
	if err != nil {
		return fmt.Errorf("publish live event: %w", err)
	}
	if _, err := tx.Exec(ctx, `SELECT pg_notify($1, s) FROM unnest($2::text[]) AS s`, notifyChannel, signals); err != nil {
		return fmt.Errorf("signal live event: %w", err)
	}
	return nil
}

func (Store) After(ctx context.Context, tx pgx.Tx, userID string, after int64) ([]Event, error) {
	rows, err := tx.Query(ctx, `
		SELECT id, type, payload, at FROM core.live_events
		WHERE user_id::text = $1 AND id > $2 ORDER BY id LIMIT $3`, userID, after, replayLimit+1)
	if err != nil {
		return nil, fmt.Errorf("read live events: %w", err)
	}
	return pgx.CollectRows(rows, scanEvent)
}

func (Store) Latest(ctx context.Context, tx pgx.Tx, userID string) (int64, error) {
	var id int64
	err := tx.QueryRow(ctx, `SELECT coalesce(max(id), 0) FROM core.live_events WHERE user_id::text = $1`, userID).Scan(&id)
	return id, err
}

func (Store) ByID(ctx context.Context, tx pgx.Tx, id int64) (Event, error) {
	rows, err := tx.Query(ctx, `SELECT id, type, payload, at FROM core.live_events WHERE id = $1`, id)
	if err != nil {
		return Event{}, fmt.Errorf("read live event: %w", err)
	}
	return pgx.CollectExactlyOneRow(rows, scanEvent)
}

func (Store) Prune(ctx context.Context, tx pgx.Tx, keep time.Duration) error {
	_, err := tx.Exec(ctx, `DELETE FROM core.live_events WHERE at < now() - $1::interval`, keep)
	return err
}

func scanEvent(row pgx.CollectableRow) (Event, error) {
	var e Event
	err := row.Scan(&e.ID, &e.Type, &e.Payload, &e.At)
	return e, err
}
