package team

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type CardStore struct{}

func (CardStore) Get(ctx context.Context, tx pgx.Tx, userID string) (Card, error) {
	cards, err := CardStore{}.ByUsers(ctx, tx, []string{userID})
	return cards[userID], err
}

func (CardStore) Put(ctx context.Context, tx pgx.Tx, userID string, c Card) error {
	if c.empty() {
		_, err := tx.Exec(ctx, `DELETE FROM team.cards WHERE user_id::text = $1`, userID)
		if err != nil {
			return fmt.Errorf("clear card: %w", err)
		}
		return nil
	}
	_, err := tx.Exec(ctx, `
		INSERT INTO team.cards (user_id, interests, skills, hours, role) VALUES ($1::uuid, $2, $3, $4, $5)
		ON CONFLICT (tenant_id, user_id) DO UPDATE
		SET interests = excluded.interests, skills = excluded.skills, hours = excluded.hours, role = excluded.role, updated_at = now()`,
		userID, c.Interests, c.Skills, c.Hours, c.Role)
	if err != nil {
		return fmt.Errorf("save card: %w", err)
	}
	return nil
}

func (CardStore) ByUsers(ctx context.Context, tx pgx.Tx, ids []string) (map[string]Card, error) {
	out := map[string]Card{}
	rows, err := tx.Query(ctx, `
		SELECT user_id::text, interests, skills, hours, role, updated_at FROM team.cards WHERE user_id::text = ANY($1)`, ids)
	if err != nil {
		return nil, fmt.Errorf("cards: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		var c Card
		if err := rows.Scan(&id, &c.Interests, &c.Skills, &c.Hours, &c.Role, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("cards: %w", err)
		}
		out[id] = c
	}
	return out, rows.Err()
}
