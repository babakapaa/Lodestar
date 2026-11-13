package study

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type UpcomingStore struct{}

func (UpcomingStore) Between(ctx context.Context, tx pgx.Tx, from, to time.Time) ([]UpcomingDefense, error) {
	rows, err := tx.Query(ctx, `
		SELECT d.assignment_id::text, d.student_id::text, s.id::text, a.title, c.title, s.starts_at
		FROM study.defenses d
		JOIN study.defense_slots s ON s.id = d.slot_id
		JOIN study.assignments a ON a.id = d.assignment_id
		JOIN study.courses c ON c.id = a.course_id
		WHERE d.status = 'planned' AND c.archived_at IS NULL
		  AND s.starts_at >= $1 AND s.starts_at < $2
		  AND d.reminded_slot IS DISTINCT FROM d.slot_id
		ORDER BY s.starts_at`, from, to)
	if err != nil {
		return nil, fmt.Errorf("upcoming defenses: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (UpcomingDefense, error) {
		var u UpcomingDefense
		err := row.Scan(&u.AssignmentID, &u.StudentID, &u.SlotID, &u.Assignment, &u.Course, &u.StartsAt)
		return u, err
	})
}

func (UpcomingStore) Mark(ctx context.Context, tx pgx.Tx, u UpcomingDefense) error {
	_, err := tx.Exec(ctx, `
		UPDATE study.defenses SET reminded_slot = slot_id
		WHERE assignment_id::text = $1 AND student_id::text = $2 AND slot_id::text = $3`,
		u.AssignmentID, u.StudentID, u.SlotID)
	if err != nil {
		return fmt.Errorf("mark defense reminder: %w", err)
	}
	return nil
}
