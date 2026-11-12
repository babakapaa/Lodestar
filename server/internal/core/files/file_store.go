package files

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

const (
	pageSize    = 50
	fileColumns = `id::text, owner_id::text, name, content_type, size, status, created_at, uploaded_at`
)

type FileStore struct{}

func (FileStore) Create(ctx context.Context, tx pgx.Tx, ownerID string, in FileInput) (File, error) {
	row := tx.QueryRow(ctx, `
		INSERT INTO core.files (owner_id, name, content_type, size) VALUES ($1, $2, $3, $4)
		RETURNING `+fileColumns, ownerID, in.Name, in.ContentType, in.Size)
	f, err := scanFile(row)
	if err != nil {
		return File{}, fmt.Errorf("save file: %w", err)
	}
	return f, nil
}

func (FileStore) ByID(ctx context.Context, tx pgx.Tx, id string, lock bool) (File, error) {
	query := `SELECT ` + fileColumns + ` FROM core.files WHERE id::text = $1`
	if lock {
		query += ` FOR UPDATE`
	}
	f, err := scanFile(tx.QueryRow(ctx, query, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return File{}, ErrNotFound
	}
	if err != nil {
		return File{}, fmt.Errorf("read file: %w", err)
	}
	return f, nil
}

func (FileStore) MarkReady(ctx context.Context, tx pgx.Tx, id string) (File, error) {
	row := tx.QueryRow(ctx, `
		UPDATE core.files SET status = 'ready', uploaded_at = now()
		WHERE id::text = $1 RETURNING `+fileColumns, id)
	f, err := scanFile(row)
	if err != nil {
		return File{}, fmt.Errorf("mark file ready: %w", err)
	}
	return f, nil
}

func (FileStore) Delete(ctx context.Context, tx pgx.Tx, id string) error {
	_, err := tx.Exec(ctx, `DELETE FROM core.files WHERE id::text = $1`, id)
	return err
}

func (FileStore) Page(ctx context.Context, tx pgx.Tx, ownerID string, before *time.Time) ([]File, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+fileColumns+` FROM core.files
		WHERE owner_id::text = $1 AND status = 'ready' AND ($2::timestamptz IS NULL OR created_at < $2)
		ORDER BY created_at DESC, id LIMIT $3`, ownerID, before, pageSize)
	if err != nil {
		return nil, fmt.Errorf("list files: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (File, error) {
		return scanFile(row)
	})
}

func (FileStore) Abandoned(ctx context.Context, tx pgx.Tx, olderThan time.Duration) ([]string, error) {
	rows, err := tx.Query(ctx, `
		DELETE FROM core.files WHERE status = 'pending' AND created_at < now() - $1::interval
		RETURNING id::text`, olderThan)
	if err != nil {
		return nil, fmt.Errorf("drop abandoned uploads: %w", err)
	}
	return pgx.CollectRows(rows, pgx.RowTo[string])
}

func scanFile(row pgx.Row) (File, error) {
	var f File
	err := row.Scan(&f.ID, &f.OwnerID, &f.Name, &f.ContentType, &f.Size, &f.Status, &f.CreatedAt, &f.UploadedAt)
	return f, err
}

func (FileStore) Owned(ctx context.Context, tx pgx.Tx, ownerID string) ([]File, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+fileColumns+` FROM core.files WHERE owner_id::text = $1 ORDER BY created_at, id`, ownerID)
	if err != nil {
		return nil, fmt.Errorf("list owned files: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (File, error) {
		return scanFile(row)
	})
}

func (FileStore) DeleteOwned(ctx context.Context, tx pgx.Tx, ownerID string) error {
	_, err := tx.Exec(ctx, `DELETE FROM core.files WHERE owner_id::text = $1`, ownerID)
	return err
}
