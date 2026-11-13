package module

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type ScoreSource interface {
	RosterSource
	ScoredWorks(ctx context.Context, tx pgx.Tx, rosterID string) ([]ScoredWork, error)
}
