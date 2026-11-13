package navigator

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/module"
)

const MetricUsers = "navigator.users"

type NavigatorStats struct{}

func (NavigatorStats) Stats(ctx context.Context, tx pgx.Tx, userIDs []string, from, to time.Time) ([]module.Metric, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}
	var users int
	err := tx.QueryRow(ctx, `
		SELECT count(DISTINCT d.owner_id)
		FROM navigator.dialogs d JOIN navigator.messages m ON m.dialog_id = d.id
		WHERE d.owner_id = ANY($1::uuid[]) AND m.role = 'student' AND m.created_at >= $2 AND m.created_at < $3`,
		userIDs, from, to).Scan(&users)
	if err != nil {
		return nil, err
	}
	return []module.Metric{{Key: MetricUsers, Value: float64(users) / float64(len(userIDs))}}, nil
}
