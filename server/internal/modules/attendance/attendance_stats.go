package attendance

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/module"
)

const MetricRate = "attendance.rate"

type AttendanceStats struct{}

func (AttendanceStats) Stats(ctx context.Context, tx pgx.Tx, userIDs []string, from, to time.Time) ([]module.Metric, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}
	var marked, attended int
	err := tx.QueryRow(ctx, `
		SELECT count(*) FILTER (WHERE m.status <> 'excused'),
		       count(*) FILTER (WHERE m.status IN ('present', 'late'))
		FROM attendance.marks m JOIN attendance.sessions s ON s.id = m.session_id
		WHERE m.user_id = ANY($1::uuid[]) AND m.status IS NOT NULL AND s.starts_at >= $2 AND s.starts_at < $3`,
		userIDs, from, to).Scan(&marked, &attended)
	if err != nil || marked == 0 {
		return nil, err
	}
	return []module.Metric{{Key: MetricRate, Value: float64(attended) / float64(marked)}}, nil
}
