package study

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	uniqueViolation     = "23505"
	checkViolation      = "23514"
	foreignKeyViolation = "23503"
	invalidText         = "22P02"
)

const courseColumns = `c.id::text, c.unit_id::text, c.title, c.code, c.description, c.archived_at, c.created_at`

type CourseStore struct{}

func (CourseStore) ByID(ctx context.Context, tx pgx.Tx, id string) (Course, error) {
	rows, err := tx.Query(ctx, `SELECT `+courseColumns+` FROM study.courses c WHERE c.id::text = $1`, id)
	if err != nil {
		return Course{}, fmt.Errorf("find course: %w", err)
	}
	c, err := pgx.CollectExactlyOneRow(rows, scanCourse)
	if errors.Is(err, pgx.ErrNoRows) {
		return Course{}, ErrCourseNotFound
	}
	return c, err
}

func (CourseStore) Mine(ctx context.Context, tx pgx.Tx, userID string, units []string, archived bool) ([]MyCourse, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+courseColumns+`, CASE WHEN t.user_id IS NULL THEN 'student' ELSE 'teacher' END
		FROM study.courses c
		LEFT JOIN study.course_teachers t ON t.course_id = c.id AND t.user_id::text = $1
		WHERE ($3 OR c.archived_at IS NULL)
		  AND (t.user_id IS NOT NULL
		       OR EXISTS (SELECT 1 FROM study.course_students s WHERE s.course_id = c.id AND s.user_id::text = $1)
		       OR EXISTS (SELECT 1 FROM study.course_units u WHERE u.course_id = c.id AND u.unit_id::text = ANY($2)))
		ORDER BY c.title, c.id`, userID, units, archived)
	if err != nil {
		return nil, fmt.Errorf("my courses: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (MyCourse, error) {
		var m MyCourse
		err := row.Scan(&m.ID, &m.UnitID, &m.Title, &m.Code, &m.Description, &m.ArchivedAt, &m.CreatedAt, &m.Role)
		return m, err
	})
}

func (CourseStore) ByUnit(ctx context.Context, tx pgx.Tx, unitID *string, archived bool) ([]Course, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+courseColumns+` FROM study.courses c
		WHERE ($1::text IS NULL OR c.unit_id::text = $1) AND ($2 OR c.archived_at IS NULL)
		ORDER BY c.title, c.id`, unitID, archived)
	if err != nil {
		return nil, fmt.Errorf("list courses: %w", err)
	}
	return pgx.CollectRows(rows, scanCourse)
}

func (CourseStore) Insert(ctx context.Context, tx pgx.Tx, in CourseInput) (Course, error) {
	rows, err := tx.Query(ctx, `
		INSERT INTO study.courses AS c (unit_id, title, code, description) VALUES ($1, $2, $3, $4)
		RETURNING `+courseColumns, in.UnitID, in.Title, in.Code, in.Description)
	if err != nil {
		return Course{}, courseError(err)
	}
	c, err := pgx.CollectExactlyOneRow(rows, scanCourse)
	return c, courseError(err)
}

func (CourseStore) Update(ctx context.Context, tx pgx.Tx, id string, in CourseInput) (Course, error) {
	rows, err := tx.Query(ctx, `
		UPDATE study.courses c SET unit_id = $2, title = $3, code = $4, description = $5, updated_at = now()
		WHERE c.id::text = $1 RETURNING `+courseColumns, id, in.UnitID, in.Title, in.Code, in.Description)
	if err != nil {
		return Course{}, courseError(err)
	}
	c, err := pgx.CollectExactlyOneRow(rows, scanCourse)
	if errors.Is(err, pgx.ErrNoRows) {
		return Course{}, ErrCourseNotFound
	}
	return c, courseError(err)
}

func (CourseStore) SetArchived(ctx context.Context, tx pgx.Tx, id string, archived bool) error {
	_, err := tx.Exec(ctx, `
		UPDATE study.courses SET archived_at = CASE WHEN $2 THEN coalesce(archived_at, now()) END, updated_at = now()
		WHERE id::text = $1`, id, archived)
	if err != nil {
		return fmt.Errorf("archive course: %w", err)
	}
	return nil
}

func scanCourse(row pgx.CollectableRow) (Course, error) {
	var c Course
	err := row.Scan(&c.ID, &c.UnitID, &c.Title, &c.Code, &c.Description, &c.ArchivedAt, &c.CreatedAt)
	return c, err
}

func courseError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case uniqueViolation:
			return ErrCodeTaken
		case checkViolation:
			return ErrInvalidCourse
		case foreignKeyViolation, invalidText:
			return ErrUnitNotFound
		}
	}
	if err != nil {
		return fmt.Errorf("save course: %w", err)
	}
	return nil
}
