package study

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type ReminderStore struct{}

func (ReminderStore) Pending(ctx context.Context, tx pgx.Tx, assignmentID string, students []string) ([]string, error) {
	return ids(ctx, tx, `
		SELECT s.id FROM unnest($2::text[]) AS s(id)
		WHERE NOT EXISTS (SELECT 1 FROM study.reminders r WHERE r.assignment_id::text = $1 AND r.student_id::text = s.id)
		  AND NOT EXISTS (SELECT 1 FROM study.submissions w WHERE w.assignment_id::text = $1 AND w.student_id::text = s.id)`,
		assignmentID, students)
}

func (ReminderStore) Mark(ctx context.Context, tx pgx.Tx, assignmentID string, students []string) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO study.reminders (assignment_id, student_id)
		SELECT $1, s.id::uuid FROM unnest($2::text[]) AS s(id) ON CONFLICT DO NOTHING`, assignmentID, students)
	if err != nil {
		return fmt.Errorf("mark reminders: %w", err)
	}
	return nil
}
