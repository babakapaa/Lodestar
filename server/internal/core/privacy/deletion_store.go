package privacy

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type DeletionStore struct{}

func (s DeletionStore) Put(ctx context.Context, tx pgx.Tx, userID string, grace time.Duration) (Deletion, error) {
	_, err := tx.Exec(ctx, `
		INSERT INTO core.deletion_requests (user_id, due_at) VALUES ($1, now() + $2::interval)
		ON CONFLICT (user_id) DO NOTHING`, userID, grace)
	if err != nil {
		return Deletion{}, fmt.Errorf("save deletion request: %w", err)
	}
	return s.Get(ctx, tx, userID)
}

func (DeletionStore) Get(ctx context.Context, tx pgx.Tx, userID string) (Deletion, error) {
	var d Deletion
	err := tx.QueryRow(ctx, `SELECT requested_at, due_at FROM core.deletion_requests WHERE user_id::text = $1`, userID).
		Scan(&d.RequestedAt, &d.DueAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return Deletion{}, ErrNoRequest
	}
	return d, err
}

func (DeletionStore) Delete(ctx context.Context, tx pgx.Tx, userID string) (bool, error) {
	tag, err := tx.Exec(ctx, `DELETE FROM core.deletion_requests WHERE user_id::text = $1`, userID)
	return tag.RowsAffected() == 1, err
}

func (DeletionStore) Due(ctx context.Context, tx pgx.Tx, now time.Time) ([]string, error) {
	rows, err := tx.Query(ctx, `
		SELECT user_id::text FROM core.deletion_requests WHERE due_at <= $1 ORDER BY due_at`, now)
	if err != nil {
		return nil, fmt.Errorf("read due deletions: %w", err)
	}
	return pgx.CollectRows(rows, pgx.RowTo[string])
}
