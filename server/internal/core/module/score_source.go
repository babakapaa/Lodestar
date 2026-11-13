package module

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type ScoreSource interface {
	Kind() string
	ScoredWorks(ctx context.Context, tx pgx.Tx, rosterID string) ([]ScoredWork, error)
}
