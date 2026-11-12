package study

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const submissionColumns = `s.id::text, s.assignment_id::text, s.student_id::text, s.client_id::text, s.body, s.composed_at, s.submitted_at`

type SubmissionStore struct{}

func (SubmissionStore) ByClient(ctx context.Context, tx pgx.Tx, assignmentID, studentID, clientID string) (Submission, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+submissionColumns+` FROM study.submissions s
		WHERE s.assignment_id::text = $1 AND s.student_id::text = $2 AND s.client_id::text = $3`, assignmentID, studentID, clientID)
	if err != nil {
		return Submission{}, fmt.Errorf("find submission: %w", err)
	}
	sub, err := pgx.CollectExactlyOneRow(rows, scanSubmission)
	if errors.Is(err, pgx.ErrNoRows) {
		return Submission{}, ErrAssignmentNotFound
	}
	return sub, err
}

func (SubmissionStore) Insert(ctx context.Context, tx pgx.Tx, assignmentID, studentID string, in SubmissionInput) (Submission, error) {
	rows, err := tx.Query(ctx, `
		INSERT INTO study.submissions AS s (assignment_id, student_id, client_id, body, composed_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING `+submissionColumns, assignmentID, studentID, in.ClientID, in.Body, in.ComposedAt)
	if err != nil {
		return Submission{}, submissionError(err)
	}
	sub, err := pgx.CollectExactlyOneRow(rows, scanSubmission)
	return sub, submissionError(err)
}

func (SubmissionStore) SetFiles(ctx context.Context, tx pgx.Tx, id string, fileIDs []string) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO study.submission_files (submission_id, file_id, position)
		SELECT $1, f.id, f.position FROM unnest($2::uuid[]) WITH ORDINALITY AS f(id, position)`, id, fileIDs)
	if err != nil {
		return fmt.Errorf("attach submission files: %w", err)
	}
	return nil
}

func (SubmissionStore) Files(ctx context.Context, tx pgx.Tx, id string) ([]string, error) {
	return ids(ctx, tx, `SELECT file_id::text FROM study.submission_files WHERE submission_id::text = $1 ORDER BY position`, id)
}

func (SubmissionStore) Attempts(ctx context.Context, tx pgx.Tx, assignmentID, studentID string) ([]Submission, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+submissionColumns+` FROM study.submissions s
		WHERE s.assignment_id::text = $1 AND s.student_id::text = $2
		ORDER BY s.submitted_at DESC, s.id`, assignmentID, studentID)
	if err != nil {
		return nil, fmt.Errorf("list attempts: %w", err)
	}
	return pgx.CollectRows(rows, scanSubmission)
}

func (SubmissionStore) Latest(ctx context.Context, tx pgx.Tx, assignmentID string) (map[string]Submission, map[string]int, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+submissionColumns+`, count(*) OVER (PARTITION BY s.student_id)
		FROM study.submissions s WHERE s.assignment_id::text = $1
		ORDER BY s.student_id, s.submitted_at DESC, s.id`, assignmentID)
	if err != nil {
		return nil, nil, fmt.Errorf("latest submissions: %w", err)
	}
	defer rows.Close()
	latest := map[string]Submission{}
	attempts := map[string]int{}
	for rows.Next() {
		var sub Submission
		var count int
		if err := rows.Scan(&sub.ID, &sub.AssignmentID, &sub.StudentID, &sub.ClientID, &sub.Body, &sub.ComposedAt, &sub.SubmittedAt, &count); err != nil {
			return nil, nil, fmt.Errorf("scan submission: %w", err)
		}
		if _, seen := latest[sub.StudentID]; !seen {
			latest[sub.StudentID] = sub
			attempts[sub.StudentID] = count
		}
	}
	return latest, attempts, rows.Err()
}

func (SubmissionStore) LastSubmitted(ctx context.Context, tx pgx.Tx, assignmentIDs []string, studentID string) (map[string]time.Time, error) {
	rows, err := tx.Query(ctx, `
		SELECT assignment_id::text, max(submitted_at) FROM study.submissions
		WHERE assignment_id::text = ANY($1) AND student_id::text = $2 GROUP BY assignment_id`, assignmentIDs, studentID)
	if err != nil {
		return nil, fmt.Errorf("last submitted: %w", err)
	}
	defer rows.Close()
	out := map[string]time.Time{}
	for rows.Next() {
		var id string
		var at time.Time
		if err := rows.Scan(&id, &at); err != nil {
			return nil, fmt.Errorf("scan last submitted: %w", err)
		}
		out[id] = at
	}
	return out, rows.Err()
}

func (SubmissionStore) Count(ctx context.Context, tx pgx.Tx, assignmentID string) (int, error) {
	var n int
	err := tx.QueryRow(ctx, `SELECT count(*) FROM study.submissions WHERE assignment_id::text = $1`, assignmentID).Scan(&n)
	if err != nil {
		return 0, fmt.Errorf("count submissions: %w", err)
	}
	return n, nil
}

func (SubmissionStore) Submitters(ctx context.Context, tx pgx.Tx, assignmentIDs []string) (map[string]int, error) {
	return countBy(ctx, tx, `
		SELECT assignment_id::text, count(DISTINCT student_id) FROM study.submissions
		WHERE assignment_id::text = ANY($1) GROUP BY assignment_id`, assignmentIDs)
}

func scanSubmission(row pgx.CollectableRow) (Submission, error) {
	var s Submission
	err := row.Scan(&s.ID, &s.AssignmentID, &s.StudentID, &s.ClientID, &s.Body, &s.ComposedAt, &s.SubmittedAt)
	return s, err
}

func submissionError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && (pgErr.Code == checkViolation || pgErr.Code == invalidText) {
		return ErrInvalidSubmission
	}
	if err != nil {
		return fmt.Errorf("save submission: %w", err)
	}
	return nil
}

func countBy(ctx context.Context, tx pgx.Tx, sql string, args ...any) (map[string]int, error) {
	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("count: %w", err)
	}
	defer rows.Close()
	out := map[string]int{}
	for rows.Next() {
		var id string
		var n int
		if err := rows.Scan(&id, &n); err != nil {
			return nil, fmt.Errorf("scan count: %w", err)
		}
		out[id] = n
	}
	return out, rows.Err()
}
