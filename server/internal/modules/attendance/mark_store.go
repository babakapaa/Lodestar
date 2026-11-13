package attendance

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

const markColumns = `m.session_id::text, m.user_id::text, m.status, m.method, m.marked_at`

type MarkStore struct{}

func (MarkStore) Expect(ctx context.Context, tx pgx.Tx, sessionID string, userIDs []string) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO attendance.marks (session_id, user_id)
		SELECT $1::uuid, u::uuid FROM unnest($2::text[]) u
		ON CONFLICT DO NOTHING`, sessionID, userIDs)
	if err != nil {
		return fmt.Errorf("expect students: %w", err)
	}
	return nil
}

func (MarkStore) OfSessions(ctx context.Context, tx pgx.Tx, ids []string) ([]markRow, error) {
	rows, err := tx.Query(ctx, `SELECT `+markColumns+` FROM attendance.marks m WHERE m.session_id::text = ANY($1)`, ids)
	if err != nil {
		return nil, fmt.Errorf("attendance marks: %w", err)
	}
	return pgx.CollectRows(rows, scanMark)
}

func (MarkStore) One(ctx context.Context, tx pgx.Tx, sessionID, userID string) (markRow, bool, error) {
	rows, err := tx.Query(ctx, `SELECT `+markColumns+` FROM attendance.marks m WHERE m.session_id::text = $1 AND m.user_id::text = $2`, sessionID, userID)
	if err != nil {
		return markRow{}, false, fmt.Errorf("attendance mark: %w", err)
	}
	m, err := pgx.CollectExactlyOneRow(rows, scanMark)
	if errors.Is(err, pgx.ErrNoRows) {
		return markRow{}, false, nil
	}
	return m, err == nil, err
}

func (MarkStore) Put(ctx context.Context, tx pgx.Tx, sessionID, userID string, status *MarkStatus, method *MarkMethod, by string, at time.Time) error {
	var stamp *time.Time
	if status != nil {
		stamp = &at
	}
	_, err := tx.Exec(ctx, `
		INSERT INTO attendance.marks (session_id, user_id, status, method, marked_by, marked_at) VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (session_id, user_id) DO UPDATE SET
			status = excluded.status, method = excluded.method, marked_by = excluded.marked_by, marked_at = excluded.marked_at`,
		sessionID, userID, status, method, by, stamp)
	if err != nil {
		return fmt.Errorf("save attendance mark: %w", err)
	}
	return nil
}

func (MarkStore) FillAbsent(ctx context.Context, tx pgx.Tx, sessionID, by string, at time.Time) error {
	_, err := tx.Exec(ctx, `
		UPDATE attendance.marks SET status = 'absent', method = 'list', marked_by = $2, marked_at = $3
		WHERE session_id::text = $1 AND status IS NULL`, sessionID, by, at)
	if err != nil {
		return fmt.Errorf("mark absent: %w", err)
	}
	return nil
}

func (MarkStore) Mine(ctx context.Context, tx pgx.Tx, userID string) ([]MyMark, error) {
	rows, err := tx.Query(ctx, `
		SELECT s.id::text, s.course_id, s.course_title, s.starts_at, m.status, m.method
		FROM attendance.marks m JOIN attendance.sessions s ON s.id = m.session_id
		WHERE m.user_id::text = $1 AND m.status IS NOT NULL ORDER BY s.starts_at DESC LIMIT 500`, userID)
	if err != nil {
		return nil, fmt.Errorf("my attendance: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (MyMark, error) {
		var m MyMark
		err := row.Scan(&m.SessionID, &m.CourseID, &m.CourseTitle, &m.StartsAt, &m.Status, &m.Method)
		return m, err
	})
}

func scanMark(row pgx.CollectableRow) (markRow, error) {
	var m markRow
	err := row.Scan(&m.SessionID, &m.UserID, &m.Status, &m.Method, &m.MarkedAt)
	return m, err
}
