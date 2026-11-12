package notes

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type SourceStore struct{}

const sourceColumns = `id::text, owner_id::text, data, file_id::text, created_at, changed_at`

func scanSource(row pgx.Row) (sourceRow, error) {
	var r sourceRow
	err := row.Scan(&r.ID, &r.OwnerID, &r.Data, &r.FileID, &r.CreatedAt, &r.ChangedAt)
	return r, err
}

func (SourceStore) Owned(ctx context.Context, tx pgx.Tx, ownerID string) ([]sourceRow, error) {
	rows, err := tx.Query(ctx, `SELECT `+sourceColumns+` FROM notes.sources
		WHERE owner_id::text = $1 ORDER BY created_at`, ownerID)
	if err != nil {
		return nil, fmt.Errorf("list sources: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (sourceRow, error) { return scanSource(row) })
}

func (SourceStore) Get(ctx context.Context, tx pgx.Tx, id string) (sourceRow, error) {
	r, err := scanSource(tx.QueryRow(ctx, `SELECT `+sourceColumns+` FROM notes.sources
		WHERE id::text = $1 FOR UPDATE`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return r, ErrSourceNotFound
	}
	if err != nil {
		return r, fmt.Errorf("get source: %w", err)
	}
	return r, nil
}

func (SourceStore) Upsert(ctx context.Context, tx pgx.Tx, id, ownerID string, data []byte, fileID *string) (sourceRow, error) {
	r, err := scanSource(tx.QueryRow(ctx, `
		INSERT INTO notes.sources (id, owner_id, data, file_id, bytes) VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET data = excluded.data, file_id = excluded.file_id,
			bytes = excluded.bytes, changed_at = now()
		RETURNING `+sourceColumns, id, ownerID, data, fileID, int64(len(data))))
	if err != nil {
		return r, fmt.Errorf("save source: %w", err)
	}
	return r, nil
}

func (SourceStore) Delete(ctx context.Context, tx pgx.Tx, id, ownerID string) error {
	tag, err := tx.Exec(ctx, `DELETE FROM notes.sources WHERE id::text = $1 AND owner_id::text = $2`, id, ownerID)
	if err != nil {
		return fmt.Errorf("delete source: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrSourceNotFound
	}
	return nil
}

func (SourceStore) FileOwner(ctx context.Context, tx pgx.Tx, fileID string) (string, error) {
	var owner string
	err := tx.QueryRow(ctx, `SELECT owner_id::text FROM core.files WHERE id::text = $1`, fileID).Scan(&owner)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrInvalidSource
	}
	if err != nil {
		return "", fmt.Errorf("source file: %w", err)
	}
	return owner, nil
}
