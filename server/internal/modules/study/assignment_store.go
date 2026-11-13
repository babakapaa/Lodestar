package study

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/babakapaa/lodestar/server/internal/core/module"
)

const assignmentColumns = `a.id::text, a.course_id::text, a.title, a.description, a.due_at, a.grading,
	a.max_points::float8, a.search_map, a.navigator_mode, a.defense, a.team, a.published_at, a.created_at`

type AssignmentStore struct{}

func (AssignmentStore) ByID(ctx context.Context, tx pgx.Tx, id string, lock bool) (Assignment, error) {
	query := `SELECT ` + assignmentColumns + ` FROM study.assignments a WHERE a.id::text = $1`
	if lock {
		query += ` FOR UPDATE`
	}
	rows, err := tx.Query(ctx, query, id)
	if err != nil {
		return Assignment{}, fmt.Errorf("find assignment: %w", err)
	}
	a, err := pgx.CollectExactlyOneRow(rows, scanAssignment)
	if errors.Is(err, pgx.ErrNoRows) {
		return Assignment{}, ErrAssignmentNotFound
	}
	return a, err
}

func (AssignmentStore) ByCourse(ctx context.Context, tx pgx.Tx, courseID string, drafts bool) ([]Assignment, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+assignmentColumns+` FROM study.assignments a
		WHERE a.course_id::text = $1 AND ($2 OR a.published_at IS NOT NULL)
		ORDER BY a.due_at NULLS LAST, a.created_at, a.id`, courseID, drafts)
	if err != nil {
		return nil, fmt.Errorf("list assignments: %w", err)
	}
	return pgx.CollectRows(rows, scanAssignment)
}

func (AssignmentStore) Insert(ctx context.Context, tx pgx.Tx, courseID, authorID string, in AssignmentInput) (Assignment, error) {
	rows, err := tx.Query(ctx, `
		INSERT INTO study.assignments AS a (course_id, title, description, due_at, grading, max_points, published_at, created_by, search_map, navigator_mode, defense, team)
		VALUES ($1, $2, $3, $4, $5, $6, CASE WHEN $7 THEN now() END, $8, $9, $10, $11, $12)
		RETURNING `+assignmentColumns, courseID, in.Title, in.Description, in.DueAt, in.Grading, in.MaxPoints, in.Published, authorID, in.SearchMap, in.NavigatorMode, in.Defense, in.Team)
	if err != nil {
		return Assignment{}, assignmentError(err)
	}
	a, err := pgx.CollectExactlyOneRow(rows, scanAssignment)
	return a, assignmentError(err)
}

func (AssignmentStore) Update(ctx context.Context, tx pgx.Tx, id string, in AssignmentInput) (Assignment, error) {
	rows, err := tx.Query(ctx, `
		UPDATE study.assignments a SET title = $2, description = $3, due_at = $4, grading = $5, max_points = $6,
		       published_at = CASE WHEN $7 THEN coalesce(a.published_at, now()) END, search_map = $8, navigator_mode = $9, defense = $10, team = $11, updated_at = now()
		WHERE a.id::text = $1 RETURNING `+assignmentColumns,
		id, in.Title, in.Description, in.DueAt, in.Grading, in.MaxPoints, in.Published, in.SearchMap, in.NavigatorMode, in.Defense, in.Team)
	if err != nil {
		return Assignment{}, assignmentError(err)
	}
	a, err := pgx.CollectExactlyOneRow(rows, scanAssignment)
	if errors.Is(err, pgx.ErrNoRows) {
		return Assignment{}, ErrAssignmentNotFound
	}
	return a, assignmentError(err)
}

func (AssignmentStore) Retitle(ctx context.Context, tx pgx.Tx, courseID, id string, w module.ExternalWork) (Assignment, error) {
	rows, err := tx.Query(ctx, `
		UPDATE study.assignments a SET title = $3, due_at = $4, max_points = $5, updated_at = now()
		WHERE a.id::text = $2 AND a.course_id::text = $1 AND a.grading = 'points'
		RETURNING `+assignmentColumns, courseID, id, w.Title, w.DueAt, w.Max)
	if err != nil {
		return Assignment{}, assignmentError(err)
	}
	a, err := pgx.CollectExactlyOneRow(rows, scanAssignment)
	if errors.Is(err, pgx.ErrNoRows) {
		return Assignment{}, ErrAssignmentNotFound
	}
	return a, assignmentError(err)
}

func (AssignmentStore) Delete(ctx context.Context, tx pgx.Tx, id string) error {
	_, err := tx.Exec(ctx, `DELETE FROM study.assignments WHERE id::text = $1`, id)
	if err != nil {
		return fmt.Errorf("delete assignment: %w", err)
	}
	return nil
}

func (AssignmentStore) Files(ctx context.Context, tx pgx.Tx, id string) ([]string, error) {
	return ids(ctx, tx, `SELECT file_id::text FROM study.assignment_files WHERE assignment_id::text = $1 ORDER BY position`, id)
}

func (AssignmentStore) SetFiles(ctx context.Context, tx pgx.Tx, id string, fileIDs []string) error {
	if _, err := tx.Exec(ctx, `DELETE FROM study.assignment_files WHERE assignment_id::text = $1`, id); err != nil {
		return fmt.Errorf("clear assignment files: %w", err)
	}
	_, err := tx.Exec(ctx, `
		INSERT INTO study.assignment_files (assignment_id, file_id, position)
		SELECT $1, f.id, f.position FROM unnest($2::uuid[]) WITH ORDINALITY AS f(id, position)`, id, fileIDs)
	if err != nil {
		return fmt.Errorf("attach assignment files: %w", err)
	}
	return nil
}

func scanAssignment(row pgx.CollectableRow) (Assignment, error) {
	var a Assignment
	err := row.Scan(&a.ID, &a.CourseID, &a.Title, &a.Description, &a.DueAt, &a.Grading, &a.MaxPoints, &a.SearchMap, &a.NavigatorMode, &a.Defense, &a.Team, &a.PublishedAt, &a.CreatedAt)
	return a, err
}

func assignmentError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == checkViolation {
		return ErrInvalidAssignment
	}
	if err != nil {
		return fmt.Errorf("save assignment: %w", err)
	}
	return nil
}

func (AssignmentStore) Due(ctx context.Context, tx pgx.Tx, courseIDs []string, from, to time.Time) ([]Assignment, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+assignmentColumns+` FROM study.assignments a
		WHERE a.course_id::text = ANY($1) AND a.published_at IS NOT NULL AND a.due_at >= $2 AND a.due_at < $3
		ORDER BY a.due_at, a.id`, courseIDs, from, to)
	if err != nil {
		return nil, fmt.Errorf("due assignments: %w", err)
	}
	return pgx.CollectRows(rows, scanAssignment)
}

func (AssignmentStore) DueBetween(ctx context.Context, tx pgx.Tx, from, to time.Time) ([]Assignment, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+assignmentColumns+` FROM study.assignments a JOIN study.courses c ON c.id = a.course_id
		WHERE a.published_at IS NOT NULL AND c.archived_at IS NULL AND a.due_at >= $1 AND a.due_at < $2
		ORDER BY a.due_at, a.id`, from, to)
	if err != nil {
		return nil, fmt.Errorf("due assignments: %w", err)
	}
	return pgx.CollectRows(rows, scanAssignment)
}
