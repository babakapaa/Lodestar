package study

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const seriesColumns = `s.id::text, s.course_id::text, s.unit_id::text, s.teacher_id::text, s.kind, s.room,
	s.starts_on, s.ends_on, s.weekday, s.start_minute, s.minutes, s.every_weeks`

type SeriesStore struct{}

func (SeriesStore) ByID(ctx context.Context, tx pgx.Tx, id string) (Series, error) {
	rows, err := tx.Query(ctx, `SELECT `+seriesColumns+` FROM study.series s WHERE s.id::text = $1`, id)
	if err != nil {
		return Series{}, fmt.Errorf("find series: %w", err)
	}
	s, err := pgx.CollectExactlyOneRow(rows, scanSeries)
	if errors.Is(err, pgx.ErrNoRows) {
		return Series{}, ErrSeriesNotFound
	}
	return s, err
}

func (SeriesStore) ByCourses(ctx context.Context, tx pgx.Tx, courseIDs []string, from, to Day) ([]Series, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+seriesColumns+` FROM study.series s
		WHERE s.course_id::text = ANY($1) AND s.starts_on <= $3 AND s.ends_on >= $2
		ORDER BY s.weekday, s.start_minute, s.id`, courseIDs, from.Time, to.Time)
	if err != nil {
		return nil, fmt.Errorf("list series: %w", err)
	}
	return pgx.CollectRows(rows, scanSeries)
}

func (SeriesStore) OfCourse(ctx context.Context, tx pgx.Tx, courseID string) ([]Series, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+seriesColumns+` FROM study.series s WHERE s.course_id::text = $1
		ORDER BY s.starts_on, s.weekday, s.start_minute, s.id`, courseID)
	if err != nil {
		return nil, fmt.Errorf("course series: %w", err)
	}
	return pgx.CollectRows(rows, scanSeries)
}

func (SeriesStore) Insert(ctx context.Context, tx pgx.Tx, courseID string, in SeriesInput) (Series, error) {
	rows, err := tx.Query(ctx, `
		INSERT INTO study.series AS s (course_id, unit_id, teacher_id, kind, room, starts_on, ends_on, weekday, start_minute, minutes, every_weeks)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING `+seriesColumns,
		courseID, in.UnitID, in.TeacherID, in.Kind, in.Room, in.StartsOn.Time, in.EndsOn.Time, in.Weekday, int(in.StartTime), in.Minutes, in.EveryWeeks)
	if err != nil {
		return Series{}, seriesError(err)
	}
	s, err := pgx.CollectExactlyOneRow(rows, scanSeries)
	return s, seriesError(err)
}

func (SeriesStore) Update(ctx context.Context, tx pgx.Tx, id string, in SeriesInput) (Series, error) {
	rows, err := tx.Query(ctx, `
		UPDATE study.series s SET unit_id = $2, teacher_id = $3, kind = $4, room = $5, starts_on = $6, ends_on = $7,
		       weekday = $8, start_minute = $9, minutes = $10, every_weeks = $11, updated_at = now()
		WHERE s.id::text = $1 RETURNING `+seriesColumns,
		id, in.UnitID, in.TeacherID, in.Kind, in.Room, in.StartsOn.Time, in.EndsOn.Time, in.Weekday, int(in.StartTime), in.Minutes, in.EveryWeeks)
	if err != nil {
		return Series{}, seriesError(err)
	}
	s, err := pgx.CollectExactlyOneRow(rows, scanSeries)
	if errors.Is(err, pgx.ErrNoRows) {
		return Series{}, ErrSeriesNotFound
	}
	return s, seriesError(err)
}

func (SeriesStore) Delete(ctx context.Context, tx pgx.Tx, id string) error {
	_, err := tx.Exec(ctx, `DELETE FROM study.series WHERE id::text = $1`, id)
	if err != nil {
		return fmt.Errorf("delete series: %w", err)
	}
	return nil
}

func scanSeries(row pgx.CollectableRow) (Series, error) {
	var s Series
	var start int
	err := row.Scan(&s.ID, &s.CourseID, &s.UnitID, &s.TeacherID, &s.Kind, &s.Room,
		&s.StartsOn.Time, &s.EndsOn.Time, &s.Weekday, &start, &s.Minutes, &s.EveryWeeks)
	s.StartTime = Clock(start)
	return s, err
}

func seriesError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case checkViolation:
			return ErrInvalidSeries
		case foreignKeyViolation, invalidText:
			return ErrInvalidSeries
		}
	}
	if err != nil {
		return fmt.Errorf("save series: %w", err)
	}
	return nil
}
