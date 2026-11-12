package study

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type PrivacySection struct{}

func (PrivacySection) Name() string {
	return "study.courses"
}

func (PrivacySection) Export(ctx context.Context, tx pgx.Tx, userID string) (any, error) {
	rows, err := tx.Query(ctx, `
		SELECT c.title, 'teacher', t.added_at FROM study.course_teachers t JOIN study.courses c ON c.id = t.course_id
		WHERE t.user_id::text = $1
		UNION ALL
		SELECT c.title, 'student', s.added_at FROM study.course_students s JOIN study.courses c ON c.id = s.course_id
		WHERE s.user_id::text = $1
		ORDER BY 3`, userID)
	if err != nil {
		return nil, fmt.Errorf("export courses: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (map[string]any, error) {
		var title, role string
		var added any
		err := row.Scan(&title, &role, &added)
		return map[string]any{"course": title, "role": role, "since": added}, err
	})
}
