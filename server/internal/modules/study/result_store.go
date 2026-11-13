package study

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/module"
)

type ResultStore struct{}

func (ResultStore) ByCourse(ctx context.Context, tx pgx.Tx, courseID string) (map[string][]module.WorkResult, error) {
	rows, err := tx.Query(ctx, `
		SELECT a.id::text, g.student_id::text,
		       CASE WHEN a.grading = 'pass' THEN CASE WHEN g.passed THEN 1 ELSE 0 END ELSE coalesce(g.points, 0) END::float8,
		       CASE WHEN a.grading = 'pass' THEN 1 ELSE coalesce(a.max_points, 0) END::float8,
		       greatest(g.graded_at, coalesce(d.held_at, g.graded_at))
		FROM study.assignments a
		JOIN study.grades g ON g.assignment_id = a.id
		LEFT JOIN study.defenses d ON d.assignment_id = a.id AND d.student_id = g.student_id
		WHERE a.course_id::text = $1 AND a.published_at IS NOT NULL AND a.grading <> 'none'
		  AND (a.defense IS NULL OR d.status = 'held' OR (a.defense = 'sample' AND NOT coalesce(d.selected, false)))
		ORDER BY a.id, g.student_id`, courseID)
	if err != nil {
		return nil, fmt.Errorf("course results: %w", err)
	}
	out := map[string][]module.WorkResult{}
	var assignmentID string
	var r module.WorkResult
	_, err = pgx.ForEachRow(rows, []any{&assignmentID, &r.UserID, &r.Score, &r.Max, &r.At}, func() error {
		out[assignmentID] = append(out[assignmentID], r)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("course results: %w", err)
	}
	return out, nil
}
