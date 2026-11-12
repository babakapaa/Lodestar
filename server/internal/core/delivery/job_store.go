package delivery

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type JobStore struct{}

func (JobStore) Claim(ctx context.Context, tx pgx.Tx, limit int) ([]job, error) {
	rows, err := tx.Query(ctx, `
		SELECT id, user_id::text, channel, type, params, link, attempts FROM core.deliveries
		WHERE next_at <= now() ORDER BY next_at, id LIMIT $1 FOR UPDATE SKIP LOCKED`, limit)
	if err != nil {
		return nil, fmt.Errorf("claim deliveries: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (job, error) {
		var j job
		err := row.Scan(&j.id, &j.userID, &j.channel, &j.item.Type, &j.item.Params, &j.item.Link, &j.attempts)
		return j, err
	})
}

func (JobStore) Done(ctx context.Context, tx pgx.Tx, id int64) error {
	_, err := tx.Exec(ctx, `DELETE FROM core.deliveries WHERE id = $1`, id)
	return err
}

func (JobStore) Retry(ctx context.Context, tx pgx.Tx, id int64, at time.Time, reason string) error {
	_, err := tx.Exec(ctx, `
		UPDATE core.deliveries SET attempts = attempts + 1, next_at = $2, last_error = left($3, 500)
		WHERE id = $1`, id, at, reason)
	return err
}

func (JobStore) Postpone(ctx context.Context, tx pgx.Tx, id int64, at time.Time) error {
	_, err := tx.Exec(ctx, `UPDATE core.deliveries SET next_at = $2 WHERE id = $1`, id, at)
	return err
}
