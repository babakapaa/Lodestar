package study

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/module"
)

const (
	MetricSubmitted = "study.submitted"
	MetricOnTime    = "study.on_time"
)

type StudyStats struct{}

func (StudyStats) Stats(ctx context.Context, tx pgx.Tx, userIDs []string, from, to time.Time) ([]module.Metric, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}
	var submitters, timed, onTime int
	err := tx.QueryRow(ctx, `
		SELECT count(DISTINCT s.student_id),
		       count(*) FILTER (WHERE a.due_at IS NOT NULL),
		       count(*) FILTER (WHERE a.due_at IS NOT NULL AND s.submitted_at <= a.due_at)
		FROM study.submissions s JOIN study.assignments a ON a.id = s.assignment_id
		WHERE s.student_id = ANY($1::uuid[]) AND s.submitted_at >= $2 AND s.submitted_at < $3`,
		userIDs, from, to).Scan(&submitters, &timed, &onTime)
	if err != nil {
		return nil, err
	}
	out := []module.Metric{{Key: MetricSubmitted, Value: float64(submitters) / float64(len(userIDs))}}
	if timed > 0 {
		out = append(out, module.Metric{Key: MetricOnTime, Value: float64(onTime) / float64(timed)})
	}
	return out, nil
}
