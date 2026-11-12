package structure

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type EventStore struct{}

func (EventStore) Record(ctx context.Context, tx pgx.Tx, unitID string, actorID *string, action string, data any) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("encode event: %w", err)
	}
	_, err = tx.Exec(ctx, `
		INSERT INTO core.unit_events (unit_id, actor_id, action, data) VALUES ($1, $2, $3, $4)`,
		unitID, actorID, action, payload)
	if err != nil {
		return fmt.Errorf("record event: %w", err)
	}
	return nil
}

func (EventStore) List(ctx context.Context, tx pgx.Tx, unitID string, before int64, limit int) ([]Event, error) {
	rows, err := tx.Query(ctx, `
		SELECT id, actor_id::text, action, data, at FROM core.unit_events
		WHERE unit_id::text = $1 AND ($2 = 0 OR id < $2)
		ORDER BY id DESC LIMIT $3`, unitID, before, limit)
	if err != nil {
		return nil, fmt.Errorf("list events: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (Event, error) {
		var e Event
		err := row.Scan(&e.ID, &e.ActorID, &e.Action, &e.Data, &e.At)
		return e, err
	})
}
