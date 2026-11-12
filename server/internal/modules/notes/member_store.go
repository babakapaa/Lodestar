package notes

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type MemberStore struct{}

func (MemberStore) List(ctx context.Context, tx pgx.Tx, noteID string) ([]Member, error) {
	rows, err := tx.Query(ctx, `
		SELECT m.user_id::text, u.display_name, u.email, m.can_edit, m.added_at
		FROM notes.members m JOIN core.users u ON u.id = m.user_id
		WHERE m.note_id = $1 ORDER BY m.added_at`, noteID)
	if err != nil {
		return nil, fmt.Errorf("list note members: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (Member, error) {
		var m Member
		err := row.Scan(&m.UserID, &m.DisplayName, &m.Email, &m.CanEdit, &m.AddedAt)
		return m, err
	})
}

func (MemberStore) Put(ctx context.Context, tx pgx.Tx, noteID, userID string, canEdit bool, addedBy string) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO notes.members (note_id, user_id, can_edit, added_by) VALUES ($1, $2, $3, $4)
		ON CONFLICT (note_id, user_id) DO UPDATE SET can_edit = excluded.can_edit`, noteID, userID, canEdit, addedBy)
	if err != nil {
		return fmt.Errorf("put note member: %w", err)
	}
	return nil
}

func (MemberStore) Remove(ctx context.Context, tx pgx.Tx, noteID, userID string) error {
	_, err := tx.Exec(ctx, `DELETE FROM notes.members WHERE note_id = $1 AND user_id::text = $2`, noteID, userID)
	if err != nil {
		return fmt.Errorf("remove note member: %w", err)
	}
	return nil
}
