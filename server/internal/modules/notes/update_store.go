package notes

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type UpdateStore struct{}

func scanUpdate(row pgx.CollectableRow) (Update, error) {
	var u Update
	err := row.Scan(&u.Seq, &u.Key, &u.AuthorID, &u.Data, &u.CreatedAt)
	return u, err
}

func (UpdateStore) ByKey(ctx context.Context, tx pgx.Tx, noteID, key string) (Update, bool, error) {
	rows, err := tx.Query(ctx, `
		SELECT seq, client_key, author_id::text, data, created_at FROM notes.updates
		WHERE note_id = $1 AND client_key = $2`, noteID, key)
	if err != nil {
		return Update{}, false, fmt.Errorf("find note update: %w", err)
	}
	u, err := pgx.CollectExactlyOneRow(rows, scanUpdate)
	if errors.Is(err, pgx.ErrNoRows) {
		return Update{}, false, nil
	}
	return u, err == nil, err
}

func (UpdateStore) Insert(ctx context.Context, tx pgx.Tx, noteID string, seq int64, key, authorID string, data []byte) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO notes.updates (note_id, seq, client_key, author_id, data) VALUES ($1, $2, $3, $4, $5)`,
		noteID, seq, key, authorID, data)
	if err != nil {
		return fmt.Errorf("insert note update: %w", err)
	}
	return nil
}

func (UpdateStore) After(ctx context.Context, tx pgx.Tx, noteID string, after int64) ([]Update, error) {
	rows, err := tx.Query(ctx, `
		SELECT seq, client_key, author_id::text, data, created_at FROM notes.updates
		WHERE note_id = $1 AND seq > $2 ORDER BY seq`, noteID, after)
	if err != nil {
		return nil, fmt.Errorf("list note updates: %w", err)
	}
	return pgx.CollectRows(rows, scanUpdate)
}

func (UpdateStore) Prune(ctx context.Context, tx pgx.Tx, noteID string, upTo int64) (int64, error) {
	var freed int64
	err := tx.QueryRow(ctx, `
		WITH gone AS (DELETE FROM notes.updates WHERE note_id = $1 AND seq <= $2 RETURNING length(data) AS size)
		SELECT coalesce(sum(size), 0) FROM gone`, noteID, upTo).Scan(&freed)
	if err != nil {
		return 0, fmt.Errorf("prune note updates: %w", err)
	}
	return freed, nil
}
