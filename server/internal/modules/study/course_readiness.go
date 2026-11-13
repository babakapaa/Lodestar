package study

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/module"
)

type CourseReadiness struct {
	settings module.Settings
}

func (r CourseReadiness) Problems(ctx context.Context, tx pgx.Tx, tenantID string) ([]module.Problem, error) {
	if on, err := r.settings.ModuleEnabled(ctx, tenantID, Name); err != nil || !on {
		return nil, err
	}
	var lonely int
	err := tx.QueryRow(ctx, `
		SELECT count(*) FROM study.courses c
		WHERE c.archived_at IS NULL AND NOT EXISTS (SELECT 1 FROM study.course_teachers t WHERE t.course_id = c.id)`).Scan(&lonely)
	if err != nil || lonely == 0 {
		return nil, err
	}
	return []module.Problem{{Section: "teaching", Code: "courses_without_teacher", Count: lonely}}, nil
}
