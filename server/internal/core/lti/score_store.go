package lti

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type ScoreStore struct{}

func (ScoreStore) Changed(ctx context.Context, tx pgx.Tx, lineItemID, userID string, score, max float64) (bool, error) {
	var same bool
	err := tx.QueryRow(ctx, `
		SELECT score = $3 AND max = $4 FROM lti.scores WHERE line_item_id::text = $1 AND user_id::text = $2`,
		lineItemID, userID, score, max).Scan(&same)
	if errors.Is(err, pgx.ErrNoRows) {
		return true, nil
	}
	if err != nil {
		return false, fmt.Errorf("lti score: %w", err)
	}
	return !same, nil
}

func (ScoreStore) Sent(ctx context.Context, tx pgx.Tx, lineItemID, userID string, score, max float64) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO lti.scores (line_item_id, user_id, score, max) VALUES ($1, $2, $3, $4)
		ON CONFLICT (line_item_id, user_id) DO UPDATE SET score = EXCLUDED.score, max = EXCLUDED.max, sent_at = now()`,
		lineItemID, userID, score, max)
	if err != nil {
		return fmt.Errorf("save lti score: %w", err)
	}
	return nil
}
