package sandbox

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/module"
)

const MetricHours = "sandbox.hours"

type SandboxStats struct{}

func (SandboxStats) Stats(ctx context.Context, tx pgx.Tx, userIDs []string, from, to time.Time) ([]module.Metric, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}
	var hours float64
	err := tx.QueryRow(ctx, `
		SELECT coalesce(-sum(l.hours), 0)::float8
		FROM sandbox.ledger l JOIN sandbox.runs r ON r.id = l.run_id
		WHERE l.kind = 'run' AND r.user_id = ANY($1::uuid[]) AND l.created_at >= $2 AND l.created_at < $3`,
		userIDs, from, to).Scan(&hours)
	if err != nil {
		return nil, err
	}
	return []module.Metric{{Key: MetricHours, Value: hours}}, nil
}
