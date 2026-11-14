package chat

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type RecordingStore struct{}

func (RecordingStore) Active(ctx context.Context, tx pgx.Tx, callID string) (Recording, bool, error) {
	return RecordingStore{}.one(ctx, tx, `r.call_id = $1 AND r.ended_at IS NULL`, callID)
}

func (RecordingStore) ByEgress(ctx context.Context, tx pgx.Tx, egressID string) (Recording, bool, error) {
	return RecordingStore{}.one(ctx, tx, `r.egress_id = $1`, egressID)
}

func (RecordingStore) Start(ctx context.Context, tx pgx.Tx, r Recording) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO chat.call_recordings (egress_id, call_id, file_id, started_by, audio_only) VALUES ($1, $2, $3, $4, $5)`,
		r.EgressID, r.CallID, r.FileID, r.StartedBy, r.AudioOnly)
	if err != nil {
		return fmt.Errorf("start recording: %w", err)
	}
	return nil
}

func (RecordingStore) End(ctx context.Context, tx pgx.Tx, egressID string) error {
	if _, err := tx.Exec(ctx, `UPDATE chat.call_recordings SET ended_at = clock_timestamp() WHERE egress_id = $1 AND ended_at IS NULL`, egressID); err != nil {
		return fmt.Errorf("end recording: %w", err)
	}
	return nil
}

func (RecordingStore) Queue(ctx context.Context, tx pgx.Tx, egressID string) error {
	if _, err := tx.Exec(ctx, `UPDATE chat.call_recordings SET transcript = 'pending' WHERE egress_id = $1 AND transcript IS NULL`, egressID); err != nil {
		return fmt.Errorf("queue transcript: %w", err)
	}
	return nil
}

func (RecordingStore) ClaimTranscript(ctx context.Context, tx pgx.Tx, retryAfter time.Duration) (Recording, int, bool, error) {
	var egressID string
	var tries int
	err := tx.QueryRow(ctx, `
		UPDATE chat.call_recordings SET transcript_at = clock_timestamp(), transcript_tries = transcript_tries + 1
		WHERE egress_id = (
			SELECT egress_id FROM chat.call_recordings
			WHERE transcript = 'pending' AND (transcript_at IS NULL OR transcript_at < clock_timestamp() - $1::interval)
			ORDER BY ended_at LIMIT 1 FOR UPDATE SKIP LOCKED)
		RETURNING egress_id, transcript_tries`, retryAfter.String()).Scan(&egressID, &tries)
	if errors.Is(err, pgx.ErrNoRows) {
		return Recording{}, 0, false, nil
	}
	if err != nil {
		return Recording{}, 0, false, fmt.Errorf("claim transcript: %w", err)
	}
	r, found, err := RecordingStore{}.ByEgress(ctx, tx, egressID)
	return r, tries, found, err
}

func (RecordingStore) Transcribed(ctx context.Context, tx pgx.Tx, egressID, state string) error {
	if _, err := tx.Exec(ctx, `UPDATE chat.call_recordings SET transcript = $2 WHERE egress_id = $1`, egressID, state); err != nil {
		return fmt.Errorf("mark transcript: %w", err)
	}
	return nil
}

func (RecordingStore) one(ctx context.Context, tx pgx.Tx, where, arg string) (Recording, bool, error) {
	var r Recording
	err := tx.QueryRow(ctx, `
		SELECT r.egress_id, r.call_id::text, c.chat_id::text, r.file_id::text, r.started_by::text
		FROM chat.call_recordings r JOIN chat.calls c ON c.id = r.call_id
		WHERE `+where, arg).Scan(&r.EgressID, &r.CallID, &r.ChatID, &r.FileID, &r.StartedBy)
	if errors.Is(err, pgx.ErrNoRows) {
		return Recording{}, false, nil
	}
	if err != nil {
		return Recording{}, false, fmt.Errorf("read recording: %w", err)
	}
	return r, true, nil
}
