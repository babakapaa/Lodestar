package module

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type StatsSource interface {
	Stats(ctx context.Context, tx pgx.Tx, userIDs []string, from, to time.Time) ([]Metric, error)
}
