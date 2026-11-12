package study

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type WorkPrivacySection struct{}

func (WorkPrivacySection) Name() string {
	return "study.work"
}

func (WorkPrivacySection) Export(ctx context.Context, tx pgx.Tx, userID string) (any, error) {
	rows, err := tx.Query(ctx, `
		SELECT a.title, s.body, s.submitted_at, g.points::float8, g.passed, g.comment
		FROM study.submissions s
		JOIN study.assignments a ON a.id = s.assignment_id
		LEFT JOIN study.grades g ON g.assignment_id = s.assignment_id AND g.student_id = s.student_id
		WHERE s.student_id::text = $1 ORDER BY s.submitted_at`, userID)
	if err != nil {
		return nil, fmt.Errorf("export work: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (map[string]any, error) {
		var title, body string
		var at time.Time
		var points *float64
		var passed *bool
		var comment *string
		err := row.Scan(&title, &body, &at, &points, &passed, &comment)
		return map[string]any{
			"assignment": title, "text": body, "submittedAt": at,
			"points": points, "passed": passed, "feedback": comment,
		}, err
	})
}
