package study

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/core/structure"
)

type CoursePilot struct {
	settings  module.Settings
	directory structure.Directory
}

func (c CoursePilot) PilotCourses(ctx context.Context, tx pgx.Tx, tenantID string) ([]module.PilotCourse, error) {
	if on, err := c.settings.ModuleEnabled(ctx, tenantID, Name); err != nil || !on {
		return nil, err
	}
	rows, err := tx.Query(ctx, `SELECT id::text, title, code FROM study.courses WHERE archived_at IS NULL ORDER BY title, id`)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (module.PilotCourse, error) {
		var p module.PilotCourse
		err := row.Scan(&p.ID, &p.Title, &p.Code)
		return p, err
	})
}

func (c CoursePilot) InPilot(ctx context.Context, tx pgx.Tx, userID string, courseIDs []string) (bool, error) {
	units, err := c.directory.UnitsAbove(ctx, tx, userID)
	if err != nil {
		return false, err
	}
	var in bool
	err = tx.QueryRow(ctx, `
		SELECT EXISTS (SELECT 1 FROM study.course_teachers WHERE user_id::text = $1 AND course_id::text = ANY($2))
			OR EXISTS (SELECT 1 FROM study.course_students WHERE user_id::text = $1 AND course_id::text = ANY($2))
			OR EXISTS (SELECT 1 FROM study.course_units WHERE course_id::text = ANY($2) AND unit_id::text = ANY($3))`,
		userID, courseIDs, units).Scan(&in)
	return in, err
}
