package study

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

const gradeColumns = `g.assignment_id::text, g.student_id::text, g.points::float8, g.passed, g.comment, g.graded_by::text, g.graded_at`

type GradeStore struct{}

func (GradeStore) Upsert(ctx context.Context, tx pgx.Tx, assignmentID, studentID, graderID string, in GradeInput) (Grade, error) {
	rows, err := tx.Query(ctx, `
		INSERT INTO study.grades AS g (assignment_id, student_id, points, passed, comment, graded_by)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (assignment_id, student_id) DO UPDATE
		SET points = EXCLUDED.points, passed = EXCLUDED.passed, comment = EXCLUDED.comment,
		    graded_by = EXCLUDED.graded_by, graded_at = now()
		RETURNING `+gradeColumns, assignmentID, studentID, in.Points, in.Passed, in.Comment, graderID)
	if err != nil {
		return Grade{}, fmt.Errorf("save grade: %w", err)
	}
	g, err := pgx.CollectExactlyOneRow(rows, scanGrade)
	if err != nil {
		return Grade{}, fmt.Errorf("save grade: %w", err)
	}
	return g, nil
}

func (GradeStore) ByAssignment(ctx context.Context, tx pgx.Tx, assignmentIDs []string, studentID *string) ([]Grade, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+gradeColumns+` FROM study.grades g
		WHERE g.assignment_id::text = ANY($1) AND ($2::text IS NULL OR g.student_id::text = $2)`, assignmentIDs, studentID)
	if err != nil {
		return nil, fmt.Errorf("list grades: %w", err)
	}
	return pgx.CollectRows(rows, scanGrade)
}

func scanGrade(row pgx.CollectableRow) (Grade, error) {
	var g Grade
	err := row.Scan(&g.AssignmentID, &g.StudentID, &g.Points, &g.Passed, &g.Comment, &g.GradedBy, &g.GradedAt)
	return g, err
}
