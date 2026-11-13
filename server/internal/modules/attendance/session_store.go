package attendance

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/module"
)

const sessionColumns = `s.id::text, s.series_id, s.lesson_date::text, s.course_id, s.course_title, s.room,
	s.starts_at, s.ends_at, s.secret, s.code, s.code_until, s.closed_at`

type SessionStore struct{}

func (SessionStore) Open(ctx context.Context, tx pgx.Tx, l module.LessonRef, secret []byte, code string, until time.Time, by string) (sessionRow, error) {
	rows, err := tx.Query(ctx, `
		INSERT INTO attendance.sessions AS s (series_id, lesson_date, course_id, course_title, room, starts_at, ends_at, secret, code, code_until, opened_by)
		VALUES ($1, $2::date, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (tenant_id, series_id, lesson_date) DO UPDATE SET
			course_title = excluded.course_title, room = excluded.room, starts_at = excluded.starts_at, ends_at = excluded.ends_at,
			code = excluded.code, code_until = excluded.code_until, closed_at = NULL
		RETURNING `+sessionColumns,
		l.SeriesID, l.Date, l.CourseID, l.CourseTitle, l.Room, l.StartsAt, l.EndsAt, secret, code, until, by)
	if err != nil {
		return sessionRow{}, fmt.Errorf("open attendance session: %w", err)
	}
	return pgx.CollectExactlyOneRow(rows, scanSession)
}

func (SessionStore) ByID(ctx context.Context, tx pgx.Tx, id string) (sessionRow, error) {
	rows, err := tx.Query(ctx, `SELECT `+sessionColumns+` FROM attendance.sessions s WHERE s.id::text = $1`, id)
	if err != nil {
		return sessionRow{}, fmt.Errorf("find attendance session: %w", err)
	}
	r, err := pgx.CollectExactlyOneRow(rows, scanSession)
	if errors.Is(err, pgx.ErrNoRows) {
		return sessionRow{}, ErrSessionNotFound
	}
	return r, err
}

func (SessionStore) ByLesson(ctx context.Context, tx pgx.Tx, seriesID, date string) (sessionRow, bool, error) {
	rows, err := tx.Query(ctx, `SELECT `+sessionColumns+` FROM attendance.sessions s WHERE s.series_id = $1 AND s.lesson_date = $2::date`, seriesID, date)
	if err != nil {
		return sessionRow{}, false, fmt.Errorf("find attendance session: %w", err)
	}
	r, err := pgx.CollectExactlyOneRow(rows, scanSession)
	if errors.Is(err, pgx.ErrNoRows) {
		return sessionRow{}, false, nil
	}
	return r, err == nil, err
}

func (SessionStore) OfCourse(ctx context.Context, tx pgx.Tx, courseID, from, to string) ([]sessionRow, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+sessionColumns+` FROM attendance.sessions s
		WHERE s.course_id = $1 AND s.lesson_date BETWEEN $2::date AND $3::date ORDER BY s.starts_at`, courseID, from, to)
	if err != nil {
		return nil, fmt.Errorf("course attendance sessions: %w", err)
	}
	return pgx.CollectRows(rows, scanSession)
}

func (SessionStore) WithCode(ctx context.Context, tx pgx.Tx, code string, now time.Time) ([]sessionRow, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+sessionColumns+` FROM attendance.sessions s
		WHERE s.closed_at IS NULL AND s.code = $1 AND s.code_until > $2`, code, now)
	if err != nil {
		return nil, fmt.Errorf("find attendance code: %w", err)
	}
	return pgx.CollectRows(rows, scanSession)
}

func (SessionStore) CodeTaken(ctx context.Context, tx pgx.Tx, code string, now time.Time) (bool, error) {
	var taken bool
	err := tx.QueryRow(ctx, `
		SELECT EXISTS (SELECT 1 FROM attendance.sessions s WHERE s.closed_at IS NULL AND s.code = $1 AND s.code_until > $2)`,
		code, now).Scan(&taken)
	if err != nil {
		return false, fmt.Errorf("check attendance code: %w", err)
	}
	return taken, nil
}

func (SessionStore) SetCode(ctx context.Context, tx pgx.Tx, id, code string, until time.Time) error {
	if _, err := tx.Exec(ctx, `UPDATE attendance.sessions SET code = $2, code_until = $3 WHERE id::text = $1`, id, code, until); err != nil {
		return fmt.Errorf("refresh attendance code: %w", err)
	}
	return nil
}

func (SessionStore) Close(ctx context.Context, tx pgx.Tx, id string, at time.Time) error {
	if _, err := tx.Exec(ctx, `UPDATE attendance.sessions SET closed_at = $2 WHERE id::text = $1 AND closed_at IS NULL`, id, at); err != nil {
		return fmt.Errorf("close attendance session: %w", err)
	}
	return nil
}

func (SessionStore) OpenFor(ctx context.Context, tx pgx.Tx, userID string, after time.Time) ([]openRow, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+sessionColumns+`, m.status FROM attendance.sessions s
		JOIN attendance.marks m ON m.session_id = s.id AND m.user_id::text = $1
		WHERE s.closed_at IS NULL AND s.ends_at > $2 ORDER BY s.starts_at`, userID, after)
	if err != nil {
		return nil, fmt.Errorf("open attendance sessions: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (openRow, error) {
		var o openRow
		err := row.Scan(&o.session.ID, &o.session.SeriesID, &o.session.Date, &o.session.CourseID, &o.session.CourseTitle, &o.session.Room,
			&o.session.StartsAt, &o.session.EndsAt, &o.session.Secret, &o.session.Code, &o.session.CodeUntil, &o.session.ClosedAt, &o.status)
		return o, err
	})
}

func scanSession(row pgx.CollectableRow) (sessionRow, error) {
	var r sessionRow
	err := row.Scan(&r.ID, &r.SeriesID, &r.Date, &r.CourseID, &r.CourseTitle, &r.Room,
		&r.StartsAt, &r.EndsAt, &r.Secret, &r.Code, &r.CodeUntil, &r.ClosedAt)
	return r, err
}
