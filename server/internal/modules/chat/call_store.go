package chat

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

const callColumns = `c.id, c.chat_id, c.started_by::text, c.started_at,
	EXISTS (SELECT 1 FROM chat.call_recordings r WHERE r.call_id = c.id AND r.ended_at IS NULL)`

type CallStore struct{}

func (CallStore) Active(ctx context.Context, tx pgx.Tx, chatID string) (Call, bool, error) {
	var c Call
	err := tx.QueryRow(ctx, `
		SELECT `+callColumns+` FROM chat.calls c WHERE c.chat_id = $1 AND c.ended_at IS NULL`, chatID).
		Scan(&c.ID, &c.ChatID, &c.StartedBy, &c.StartedAt, &c.Recording)
	if errors.Is(err, pgx.ErrNoRows) {
		return Call{}, false, nil
	}
	if err != nil {
		return Call{}, false, fmt.Errorf("read call: %w", err)
	}
	c.People, err = CallStore{}.People(ctx, tx, c.ID)
	return c, err == nil, err
}

func (CallStore) ByID(ctx context.Context, tx pgx.Tx, id string) (Call, bool, error) {
	var c Call
	err := tx.QueryRow(ctx, `
		SELECT `+callColumns+` FROM chat.calls c WHERE c.id = $1 AND c.ended_at IS NULL`, id).
		Scan(&c.ID, &c.ChatID, &c.StartedBy, &c.StartedAt, &c.Recording)
	if errors.Is(err, pgx.ErrNoRows) {
		return Call{}, false, nil
	}
	if err != nil {
		return Call{}, false, fmt.Errorf("read call: %w", err)
	}
	c.People, err = CallStore{}.People(ctx, tx, c.ID)
	return c, err == nil, err
}

func (CallStore) Start(ctx context.Context, tx pgx.Tx, chatID, userID string) (string, error) {
	var id string
	err := tx.QueryRow(ctx, `
		INSERT INTO chat.calls (chat_id, started_by) VALUES ($1, $2) RETURNING id`, chatID, userID).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("start call: %w", err)
	}
	return id, nil
}

func (CallStore) Join(ctx context.Context, tx pgx.Tx, callID, userID string) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO chat.call_members (call_id, user_id) VALUES ($1, $2)
		ON CONFLICT (call_id, user_id) DO UPDATE SET joined_at = clock_timestamp(), left_at = NULL`, callID, userID)
	if err != nil {
		return fmt.Errorf("join call: %w", err)
	}
	return nil
}

func (CallStore) Leave(ctx context.Context, tx pgx.Tx, callID, userID string) error {
	_, err := tx.Exec(ctx, `
		UPDATE chat.call_members SET left_at = clock_timestamp() WHERE call_id = $1 AND user_id = $2 AND left_at IS NULL`, callID, userID)
	if err != nil {
		return fmt.Errorf("leave call: %w", err)
	}
	return nil
}

func (CallStore) People(ctx context.Context, tx pgx.Tx, callID string) ([]string, error) {
	rows, err := tx.Query(ctx, `
		SELECT user_id::text FROM chat.call_members WHERE call_id = $1 AND left_at IS NULL ORDER BY joined_at`, callID)
	if err != nil {
		return nil, fmt.Errorf("read call people: %w", err)
	}
	people, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return nil, fmt.Errorf("read call people: %w", err)
	}
	return people, nil
}

func (CallStore) End(ctx context.Context, tx pgx.Tx, callID string) error {
	if _, err := tx.Exec(ctx, `UPDATE chat.calls SET ended_at = clock_timestamp() WHERE id = $1 AND ended_at IS NULL`, callID); err != nil {
		return fmt.Errorf("end call: %w", err)
	}
	if _, err := tx.Exec(ctx, `UPDATE chat.call_members SET left_at = clock_timestamp() WHERE call_id = $1 AND left_at IS NULL`, callID); err != nil {
		return fmt.Errorf("end call: %w", err)
	}
	return nil
}
