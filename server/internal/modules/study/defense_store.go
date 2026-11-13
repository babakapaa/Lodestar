package study

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

const defenseColumns = `d.assignment_id::text, d.student_id::text, d.slot_id::text, d.selected, d.status,
	d.questions, d.summary, d.examiner_id::text, d.held_at, d.updated_at`

type DefenseStore struct{}

func (DefenseStore) Get(ctx context.Context, tx pgx.Tx, assignmentID, studentID string) (Defense, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+defenseColumns+` FROM study.defenses d
		WHERE d.assignment_id::text = $1 AND d.student_id::text = $2 FOR UPDATE`, assignmentID, studentID)
	if err != nil {
		return Defense{}, fmt.Errorf("find defense: %w", err)
	}
	d, err := pgx.CollectExactlyOneRow(rows, scanDefense)
	if errors.Is(err, pgx.ErrNoRows) {
		return Defense{AssignmentID: assignmentID, StudentID: studentID, Status: DefensePlanned, Questions: []DefenseQuestion{}}, nil
	}
	return d, err
}

func (DefenseStore) ByAssignment(ctx context.Context, tx pgx.Tx, assignmentID string) ([]Defense, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+defenseColumns+` FROM study.defenses d WHERE d.assignment_id::text = $1`, assignmentID)
	if err != nil {
		return nil, fmt.Errorf("list defenses: %w", err)
	}
	return pgx.CollectRows(rows, scanDefense)
}

func (DefenseStore) InSlot(ctx context.Context, tx pgx.Tx, slotID string) ([]string, error) {
	return ids(ctx, tx, `SELECT student_id::text FROM study.defenses WHERE slot_id::text = $1`, slotID)
}

func (s DefenseStore) Book(ctx context.Context, tx pgx.Tx, assignmentID, studentID string, slotID *string) (Defense, error) {
	return s.write(ctx, tx, `slot_id = $3::uuid`, assignmentID, studentID, slotID)
}

func (s DefenseStore) Select(ctx context.Context, tx pgx.Tx, assignmentID, studentID string, selected bool) (Defense, error) {
	return s.write(ctx, tx, `selected = $3`, assignmentID, studentID, selected)
}

func (s DefenseStore) Save(ctx context.Context, tx pgx.Tx, assignmentID, studentID string, in DefenseInput) (Defense, error) {
	questions, err := json.Marshal(in.Questions)
	if err != nil {
		return Defense{}, err
	}
	return s.write(ctx, tx, `questions = $3, summary = $4`, assignmentID, studentID, questions, in.Summary)
}

func (s DefenseStore) Finish(ctx context.Context, tx pgx.Tx, assignmentID, studentID, examinerID string, status DefenseStatus) (Defense, error) {
	return s.write(ctx, tx, `status = $3, examiner_id = $4, held_at = now()`, assignmentID, studentID, status, examinerID)
}

func (DefenseStore) write(ctx context.Context, tx pgx.Tx, set, assignmentID, studentID string, args ...any) (Defense, error) {
	if _, err := tx.Exec(ctx, `
		INSERT INTO study.defenses (assignment_id, student_id) VALUES ($1, $2)
		ON CONFLICT (assignment_id, student_id) DO NOTHING`, assignmentID, studentID); err != nil {
		return Defense{}, fmt.Errorf("start defense: %w", err)
	}
	rows, err := tx.Query(ctx, `
		UPDATE study.defenses d SET `+set+`, updated_at = now()
		WHERE d.assignment_id::text = $1 AND d.student_id::text = $2
		RETURNING `+defenseColumns, append([]any{assignmentID, studentID}, args...)...)
	if err != nil {
		return Defense{}, fmt.Errorf("save defense: %w", err)
	}
	return pgx.CollectExactlyOneRow(rows, scanDefense)
}

func scanDefense(row pgx.CollectableRow) (Defense, error) {
	var d Defense
	var questions []byte
	err := row.Scan(&d.AssignmentID, &d.StudentID, &d.SlotID, &d.Selected, &d.Status, &questions, &d.Summary, &d.ExaminerID, &d.HeldAt, &d.UpdatedAt)
	if err != nil {
		return d, err
	}
	if err := json.Unmarshal(questions, &d.Questions); err != nil {
		return d, fmt.Errorf("read defense questions: %w", err)
	}
	return d, nil
}
