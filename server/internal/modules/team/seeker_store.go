package team

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type SeekerStore struct{}

func (SeekerStore) Put(ctx context.Context, tx pgx.Tx, kind, scopeID, userID, note string) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO team.seekers (scope_kind, scope_id, user_id, note) VALUES ($1, $2, $3::uuid, $4)
		ON CONFLICT (tenant_id, scope_kind, scope_id, user_id) DO UPDATE SET note = excluded.note`, kind, scopeID, userID, note)
	if err != nil {
		return fmt.Errorf("seek team: %w", err)
	}
	return nil
}

func (SeekerStore) Delete(ctx context.Context, tx pgx.Tx, kind, scopeID, userID string) error {
	_, err := tx.Exec(ctx, `DELETE FROM team.seekers WHERE scope_kind = $1 AND scope_id = $2 AND user_id::text = $3`, kind, scopeID, userID)
	if err != nil {
		return fmt.Errorf("stop seeking: %w", err)
	}
	return nil
}

func (SeekerStore) ByScope(ctx context.Context, tx pgx.Tx, kind, scopeID string) ([]Seeker, error) {
	rows, err := tx.Query(ctx, `
		SELECT u.id::text, u.display_name, u.email, s.note FROM team.seekers s
		JOIN core.users u ON u.id = s.user_id
		WHERE s.scope_kind = $1 AND s.scope_id = $2 ORDER BY s.created_at`, kind, scopeID)
	if err != nil {
		return nil, fmt.Errorf("seekers: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (Seeker, error) {
		var s Seeker
		err := row.Scan(&s.UserID, &s.DisplayName, &s.Email, &s.Note)
		return s, err
	})
}
