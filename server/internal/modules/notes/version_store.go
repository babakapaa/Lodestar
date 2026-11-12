package notes

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type VersionStore struct{}

type versionRow struct {
	Version
	Data []byte
	Body []byte
}

func scanVersion(row pgx.CollectableRow) (versionRow, error) {
	var v versionRow
	err := row.Scan(&v.ID, &v.Seq, &v.AuthorID, &v.Bytes, &v.CreatedAt, &v.Data, &v.Body)
	return v, err
}

const versionColumns = `id::text, seq, author_id::text, bytes, created_at`

func (VersionStore) Insert(ctx context.Context, tx pgx.Tx, noteID string, seq int64, authorID string, data, body []byte) (Version, error) {
	rows, err := tx.Query(ctx, `
		INSERT INTO notes.versions (note_id, seq, author_id, data, body, bytes)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING `+versionColumns+`, data, body`, noteID, seq, authorID, data, body, len(data)+len(body))
	if err != nil {
		return Version{}, fmt.Errorf("insert note version: %w", err)
	}
	v, err := pgx.CollectExactlyOneRow(rows, scanVersion)
	return v.Version, err
}

func (VersionStore) Latest(ctx context.Context, tx pgx.Tx, noteID string) (versionRow, bool, error) {
	rows, err := tx.Query(ctx, `SELECT `+versionColumns+`, data, body FROM notes.versions
		WHERE note_id = $1 ORDER BY seq DESC LIMIT 1`, noteID)
	if err != nil {
		return versionRow{}, false, fmt.Errorf("latest note version: %w", err)
	}
	v, err := pgx.CollectExactlyOneRow(rows, scanVersion)
	if errors.Is(err, pgx.ErrNoRows) {
		return versionRow{}, false, nil
	}
	return v, err == nil, err
}

func (VersionStore) List(ctx context.Context, tx pgx.Tx, noteID string) ([]Version, error) {
	rows, err := tx.Query(ctx, `SELECT `+versionColumns+` FROM notes.versions WHERE note_id = $1 ORDER BY seq DESC`, noteID)
	if err != nil {
		return nil, fmt.Errorf("list note versions: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (Version, error) {
		var v Version
		err := row.Scan(&v.ID, &v.Seq, &v.AuthorID, &v.Bytes, &v.CreatedAt)
		return v, err
	})
}

func (VersionStore) Get(ctx context.Context, tx pgx.Tx, noteID, id string) (versionRow, error) {
	rows, err := tx.Query(ctx, `SELECT `+versionColumns+`, data, body FROM notes.versions
		WHERE note_id = $1 AND id::text = $2`, noteID, id)
	if err != nil {
		return versionRow{}, fmt.Errorf("get note version: %w", err)
	}
	v, err := pgx.CollectExactlyOneRow(rows, scanVersion)
	if errors.Is(err, pgx.ErrNoRows) {
		return versionRow{}, ErrVersionNotFound
	}
	return v, err
}
