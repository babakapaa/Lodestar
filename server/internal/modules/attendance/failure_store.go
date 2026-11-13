package attendance

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type FailureStore struct{}

func (FailureStore) Count(ctx context.Context, tx pgx.Tx, userID string, since time.Time) (int, error) {
	var n int
	err := tx.QueryRow(ctx, `SELECT count(*) FROM attendance.failures WHERE user_id::text = $1 AND failed_at > $2`, userID, since).Scan(&n)
	if err != nil {
		return 0, fmt.Errorf("count attendance failures: %w", err)
	}
	return n, nil
}

func (FailureStore) Add(ctx context.Context, tx pgx.Tx, userID string, now time.Time) error {
	if _, err := tx.Exec(ctx, `DELETE FROM attendance.failures WHERE failed_at < $1`, now.Add(-24*time.Hour)); err != nil {
		return fmt.Errorf("prune attendance failures: %w", err)
	}
	if _, err := tx.Exec(ctx, `INSERT INTO attendance.failures (user_id, failed_at) VALUES ($1, $2)`, userID, now); err != nil {
		return fmt.Errorf("record attendance failure: %w", err)
	}
	return nil
}
