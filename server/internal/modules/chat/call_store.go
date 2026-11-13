package chat

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

const callColumns = `c.id, c.chat_id, c.started_by::text, c.started_at,
	EXISTS (SELECT 1 FROM chat.call_recordings r WHERE r.call_id = c.id AND r.ended_at IS NULL),
	coalesce((SELECT r.audio_only FROM chat.call_recordings r WHERE r.call_id = c.id AND r.ended_at IS NULL), false),
	c.consent_requested_by::text, c.consent_audio_only`

type CallStore struct{}

func (s CallStore) Active(ctx context.Context, tx pgx.Tx, chatID string) (Call, bool, error) {
	return s.one(ctx, tx, `c.chat_id = $1`, chatID)
}

func (s CallStore) ByID(ctx context.Context, tx pgx.Tx, id string) (Call, bool, error) {
	return s.one(ctx, tx, `c.id = $1`, id)
}

func (s CallStore) one(ctx context.Context, tx pgx.Tx, where, arg string) (Call, bool, error) {
	var c Call
	var requestedBy *string
	var audioOnly bool
	err := tx.QueryRow(ctx, `
		SELECT `+callColumns+` FROM chat.calls c WHERE `+where+` AND c.ended_at IS NULL`, arg).
		Scan(&c.ID, &c.ChatID, &c.StartedBy, &c.StartedAt, &c.Recording, &c.AudioOnly, &requestedBy, &audioOnly)
	if errors.Is(err, pgx.ErrNoRows) {
		return Call{}, false, nil
	}
	if err != nil {
		return Call{}, false, fmt.Errorf("read call: %w", err)
	}
	if c.People, err = s.People(ctx, tx, c.ID); err != nil {
		return Call{}, false, err
	}
	if requestedBy != nil {
		waiting, err := s.Waiting(ctx, tx, c.ID)
		if err != nil {
			return Call{}, false, err
		}
		c.Consent = &ConsentRequest{RequestedBy: *requestedBy, AudioOnly: audioOnly, Waiting: waiting}
	}
	return c, true, nil
}

func (CallStore) Waiting(ctx context.Context, tx pgx.Tx, callID string) ([]string, error) {
	rows, err := tx.Query(ctx, `
		SELECT user_id::text FROM chat.call_members WHERE call_id = $1 AND left_at IS NULL AND NOT consent ORDER BY joined_at`, callID)
	if err != nil {
		return nil, fmt.Errorf("read consent: %w", err)
	}
	waiting, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return nil, fmt.Errorf("read consent: %w", err)
	}
	return waiting, nil
}

func (CallStore) Agree(ctx context.Context, tx pgx.Tx, callID, userID string) error {
	if _, err := tx.Exec(ctx, `UPDATE chat.call_members SET consent = true WHERE call_id = $1 AND user_id = $2`, callID, userID); err != nil {
		return fmt.Errorf("agree to recording: %w", err)
	}
	return nil
}

func (CallStore) RequestConsent(ctx context.Context, tx pgx.Tx, callID, userID string, audioOnly bool) error {
	_, err := tx.Exec(ctx, `UPDATE chat.calls SET consent_requested_by = $2, consent_audio_only = $3 WHERE id = $1`, callID, userID, audioOnly)
	if err != nil {
		return fmt.Errorf("request consent: %w", err)
	}
	return nil
}

func (CallStore) ClearConsent(ctx context.Context, tx pgx.Tx, callID string) error {
	if _, err := tx.Exec(ctx, `UPDATE chat.calls SET consent_requested_by = NULL, consent_audio_only = false WHERE id = $1`, callID); err != nil {
		return fmt.Errorf("clear consent: %w", err)
	}
	return nil
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
