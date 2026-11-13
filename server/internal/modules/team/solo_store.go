package team

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type SoloStore struct{}

func (SoloStore) Put(ctx context.Context, tx pgx.Tx, kind, scopeID, userID, decidedBy string) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO team.solos (scope_kind, scope_id, user_id, decided_by) VALUES ($1, $2, $3::uuid, $4::uuid)
		ON CONFLICT DO NOTHING`, kind, scopeID, userID, decidedBy)
	if err != nil {
		return fmt.Errorf("separate student: %w", err)
	}
	return nil
}

func (SoloStore) Delete(ctx context.Context, tx pgx.Tx, kind, scopeID, userID string) error {
	_, err := tx.Exec(ctx, `DELETE FROM team.solos WHERE scope_kind = $1 AND scope_id = $2 AND user_id::text = $3`, kind, scopeID, userID)
	if err != nil {
		return fmt.Errorf("return student: %w", err)
	}
	return nil
}

func (SoloStore) ByScope(ctx context.Context, tx pgx.Tx, kind, scopeID string) (map[string]bool, error) {
	rows, err := tx.Query(ctx, `SELECT user_id::text FROM team.solos WHERE scope_kind = $1 AND scope_id = $2`, kind, scopeID)
	if err != nil {
		return nil, fmt.Errorf("solos: %w", err)
	}
	ids, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return nil, fmt.Errorf("solos: %w", err)
	}
	out := make(map[string]bool, len(ids))
	for _, id := range ids {
		out[id] = true
	}
	return out, nil
}
