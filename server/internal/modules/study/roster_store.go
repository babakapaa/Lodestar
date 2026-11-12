package study

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type RosterStore struct{}

func (RosterStore) Teachers(ctx context.Context, tx pgx.Tx, courseID string) ([]string, error) {
	return ids(ctx, tx, `SELECT user_id::text FROM study.course_teachers WHERE course_id::text = $1`, courseID)
}

func (RosterStore) Units(ctx context.Context, tx pgx.Tx, courseID string) ([]string, error) {
	return ids(ctx, tx, `SELECT unit_id::text FROM study.course_units WHERE course_id::text = $1`, courseID)
}

func (RosterStore) Students(ctx context.Context, tx pgx.Tx, courseID string) ([]string, error) {
	return ids(ctx, tx, `SELECT user_id::text FROM study.course_students WHERE course_id::text = $1`, courseID)
}

func (RosterStore) Teaches(ctx context.Context, tx pgx.Tx, courseID, userID string) (bool, error) {
	return exists(ctx, tx, `SELECT 1 FROM study.course_teachers WHERE course_id::text = $1 AND user_id::text = $2`, courseID, userID)
}

func (RosterStore) Learns(ctx context.Context, tx pgx.Tx, courseID, userID string, units []string) (bool, error) {
	return exists(ctx, tx, `
		SELECT 1 FROM study.course_students WHERE course_id::text = $1 AND user_id::text = $2
		UNION ALL
		SELECT 1 FROM study.course_units WHERE course_id::text = $1 AND unit_id::text = ANY($3)`, courseID, userID, units)
}

func (RosterStore) Add(ctx context.Context, tx pgx.Tx, list List, courseID, memberID string) (bool, error) {
	tag, err := tx.Exec(ctx, `INSERT INTO study.`+list.table()+` (course_id, `+list.column()+`) VALUES ($1, $2) ON CONFLICT DO NOTHING`, courseID, memberID)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && (pgErr.Code == foreignKeyViolation || pgErr.Code == invalidText) {
		return false, list.missing()
	}
	if err != nil {
		return false, fmt.Errorf("add to %s: %w", list.table(), err)
	}
	return tag.RowsAffected() == 1, nil
}

func (RosterStore) Remove(ctx context.Context, tx pgx.Tx, list List, courseID, memberID string) (bool, error) {
	tag, err := tx.Exec(ctx, `DELETE FROM study.`+list.table()+` WHERE course_id::text = $1 AND `+list.column()+`::text = $2`, courseID, memberID)
	if err != nil {
		return false, fmt.Errorf("remove from %s: %w", list.table(), err)
	}
	return tag.RowsAffected() == 1, nil
}

func ids(ctx context.Context, tx pgx.Tx, sql string, args ...any) ([]string, error) {
	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("list ids: %w", err)
	}
	return pgx.CollectRows(rows, pgx.RowTo[string])
}

func exists(ctx context.Context, tx pgx.Tx, sql string, args ...any) (bool, error) {
	var found bool
	err := tx.QueryRow(ctx, `SELECT EXISTS (`+sql+`)`, args...).Scan(&found)
	if err != nil {
		return false, fmt.Errorf("check membership: %w", err)
	}
	return found, nil
}
