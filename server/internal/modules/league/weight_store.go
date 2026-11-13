package league

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type WeightStore struct{}

func (WeightStore) All(ctx context.Context, tx pgx.Tx, leagueID string) (Weights, error) {
	rows, err := tx.Query(ctx, `SELECT work_kind || '/' || work_id, weight FROM league.weights WHERE league_id::text = $1`, leagueID)
	if err != nil {
		return nil, fmt.Errorf("load weights: %w", err)
	}
	out := Weights{}
	var key string
	var weight int
	_, err = pgx.ForEachRow(rows, []any{&key, &weight}, func() error {
		out[key] = weight
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("load weights: %w", err)
	}
	return out, nil
}

func (WeightStore) Replace(ctx context.Context, tx pgx.Tx, leagueID string, weights []WeightInput) error {
	if _, err := tx.Exec(ctx, `DELETE FROM league.weights WHERE league_id::text = $1`, leagueID); err != nil {
		return fmt.Errorf("clear weights: %w", err)
	}
	for _, w := range weights {
		_, err := tx.Exec(ctx, `INSERT INTO league.weights (league_id, work_kind, work_id, weight) VALUES ($1, $2, $3, $4)`,
			leagueID, w.Kind, w.ID, w.Weight)
		if err != nil {
			return fmt.Errorf("save weight: %w", err)
		}
	}
	return nil
}
