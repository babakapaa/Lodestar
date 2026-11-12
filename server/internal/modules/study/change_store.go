package study

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const changeColumns = `c.series_id::text, c.on_date, c.status, c.starts_at, c.minutes, c.room, c.teacher_id::text, c.note`

type ChangeStore struct{}

func (ChangeStore) InRange(ctx context.Context, tx pgx.Tx, seriesIDs []string, from, to Day) ([]Change, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+changeColumns+` FROM study.series_changes c
		WHERE c.series_id::text = ANY($1) AND c.on_date BETWEEN $2 AND $3`, seriesIDs, from.Time, to.Time)
	if err != nil {
		return nil, fmt.Errorf("list changes: %w", err)
	}
	return pgx.CollectRows(rows, scanChange)
}

func (ChangeStore) Put(ctx context.Context, tx pgx.Tx, seriesID string, date Day, in ChangeInput) (Change, error) {
	rows, err := tx.Query(ctx, `
		INSERT INTO study.series_changes AS c (series_id, on_date, status, starts_at, minutes, room, teacher_id, note)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (series_id, on_date) DO UPDATE
		SET status = EXCLUDED.status, starts_at = EXCLUDED.starts_at, minutes = EXCLUDED.minutes, room = EXCLUDED.room,
		    teacher_id = EXCLUDED.teacher_id, note = EXCLUDED.note, changed_at = now()
		RETURNING `+changeColumns, seriesID, date.Time, in.Status, in.StartsAt, in.Minutes, in.Room, in.TeacherID, in.Note)
	if err != nil {
		return Change{}, changeError(err)
	}
	c, err := pgx.CollectExactlyOneRow(rows, scanChange)
	return c, changeError(err)
}

func (ChangeStore) Delete(ctx context.Context, tx pgx.Tx, seriesID string, date Day) (bool, error) {
	tag, err := tx.Exec(ctx, `DELETE FROM study.series_changes WHERE series_id::text = $1 AND on_date = $2`, seriesID, date.Time)
	if err != nil {
		return false, fmt.Errorf("revert change: %w", err)
	}
	return tag.RowsAffected() == 1, nil
}

func scanChange(row pgx.CollectableRow) (Change, error) {
	var c Change
	err := row.Scan(&c.SeriesID, &c.Date.Time, &c.Status, &c.StartsAt, &c.Minutes, &c.Room, &c.TeacherID, &c.Note)
	return c, err
}

func changeError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && (pgErr.Code == checkViolation || pgErr.Code == foreignKeyViolation || pgErr.Code == invalidText) {
		return ErrInvalidChange
	}
	if err != nil {
		return fmt.Errorf("save change: %w", err)
	}
	return nil
}
