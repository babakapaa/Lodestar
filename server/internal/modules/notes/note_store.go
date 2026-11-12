package notes

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type NoteStore struct{}

const noteColumns = `n.id::text, n.owner_id::text,
	CASE WHEN n.owner_id::text = $1 THEN 'owner' WHEN m.can_edit THEN 'editor' ELSE 'reader' END,
	n.meta, n.seq, n.base_seq, n.bytes, n.created_at, n.changed_at, n.deleted_at`

const noteJoin = `FROM notes.notes n LEFT JOIN notes.members m ON m.note_id = n.id AND m.user_id::text = $1`

func scanNote(row pgx.CollectableRow) (noteRow, error) {
	var n noteRow
	err := row.Scan(&n.ID, &n.OwnerID, &n.Role, &n.Meta, &n.Seq, &n.BaseSeq, &n.Bytes, &n.CreatedAt, &n.ChangedAt, &n.DeletedAt)
	return n, err
}

func (NoteStore) Insert(ctx context.Context, tx pgx.Tx, id, ownerID string, meta []byte) (bool, error) {
	tag, err := tx.Exec(ctx, `
		INSERT INTO notes.notes (id, owner_id, meta, bytes) VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO NOTHING`, id, ownerID, meta, len(meta))
	if err != nil {
		return false, fmt.Errorf("insert note: %w", err)
	}
	return tag.RowsAffected() == 1, nil
}

func (NoteStore) Visible(ctx context.Context, tx pgx.Tx, userID string) ([]noteRow, error) {
	rows, err := tx.Query(ctx, `SELECT `+noteColumns+` `+noteJoin+`
		WHERE n.owner_id::text = $1 OR m.user_id IS NOT NULL
		ORDER BY n.created_at`, userID)
	if err != nil {
		return nil, fmt.Errorf("list notes: %w", err)
	}
	return pgx.CollectRows(rows, scanNote)
}

func (NoteStore) Get(ctx context.Context, tx pgx.Tx, userID, id string, lock bool) (noteRow, error) {
	query := `SELECT ` + noteColumns + ` ` + noteJoin + `
		WHERE n.id::text = $2 AND (n.owner_id::text = $1 OR m.user_id IS NOT NULL)`
	if lock {
		query += ` FOR UPDATE OF n`
	}
	rows, err := tx.Query(ctx, query, userID, id)
	if err != nil {
		return noteRow{}, fmt.Errorf("get note: %w", err)
	}
	n, err := pgx.CollectExactlyOneRow(rows, scanNote)
	if errors.Is(err, pgx.ErrNoRows) {
		return noteRow{}, ErrNoteNotFound
	}
	return n, err
}

func (NoteStore) SetMeta(ctx context.Context, tx pgx.Tx, id string, meta []byte, bytesDelta int64) error {
	_, err := tx.Exec(ctx, `UPDATE notes.notes SET meta = $2, bytes = bytes + $3, changed_at = now() WHERE id = $1`,
		id, meta, bytesDelta)
	if err != nil {
		return fmt.Errorf("set note meta: %w", err)
	}
	return nil
}

func (NoteStore) SetDeleted(ctx context.Context, tx pgx.Tx, id string, deleted bool) error {
	_, err := tx.Exec(ctx, `
		UPDATE notes.notes SET deleted_at = CASE WHEN $2 THEN now() END, changed_at = now() WHERE id = $1`, id, deleted)
	if err != nil {
		return fmt.Errorf("delete note: %w", err)
	}
	return nil
}

func (NoteStore) Advance(ctx context.Context, tx pgx.Tx, id string, bytes int64) (int64, error) {
	var seq int64
	err := tx.QueryRow(ctx, `
		UPDATE notes.notes SET seq = seq + 1, bytes = bytes + $2, changed_at = now() WHERE id = $1 RETURNING seq`,
		id, bytes).Scan(&seq)
	if err != nil {
		return 0, fmt.Errorf("advance note: %w", err)
	}
	return seq, nil
}

func (NoteStore) Rebase(ctx context.Context, tx pgx.Tx, id string, baseSeq, bytesDelta int64) error {
	_, err := tx.Exec(ctx, `UPDATE notes.notes SET base_seq = $2, bytes = bytes + $3 WHERE id = $1`, id, baseSeq, bytesDelta)
	if err != nil {
		return fmt.Errorf("rebase note: %w", err)
	}
	return nil
}

func (NoteStore) OwnerBytes(ctx context.Context, tx pgx.Tx, ownerID string) (int64, error) {
	var total int64
	err := tx.QueryRow(ctx, `SELECT (SELECT coalesce(sum(bytes), 0) FROM notes.notes WHERE owner_id::text = $1)
		+ (SELECT coalesce(sum(bytes), 0) FROM notes.sources WHERE owner_id::text = $1)`, ownerID).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("notes usage: %w", err)
	}
	return total, nil
}

func (NoteStore) Purge(ctx context.Context, tx pgx.Tx, before time.Time) (int64, error) {
	tag, err := tx.Exec(ctx, `DELETE FROM notes.notes WHERE deleted_at < $1`, before)
	if err != nil {
		return 0, fmt.Errorf("purge notes: %w", err)
	}
	return tag.RowsAffected(), nil
}

func (NoteStore) Audience(ctx context.Context, tx pgx.Tx, id string) ([]string, error) {
	rows, err := tx.Query(ctx, `
		SELECT owner_id::text FROM notes.notes WHERE id = $1
		UNION SELECT user_id::text FROM notes.members WHERE note_id = $1`, id)
	if err != nil {
		return nil, fmt.Errorf("note audience: %w", err)
	}
	return pgx.CollectRows(rows, pgx.RowTo[string])
}
