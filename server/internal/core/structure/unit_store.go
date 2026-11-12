package structure

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

const unitColumns = `id::text, parent_id::text, kind, name, code, archived_at, created_at`

type UnitStore struct{}

func (UnitStore) Children(ctx context.Context, tx pgx.Tx, parentID *string, archived bool) ([]Unit, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+unitColumns+` FROM core.units
		WHERE parent_id::text IS NOT DISTINCT FROM $1 AND ($2 OR archived_at IS NULL)
		ORDER BY name`, parentID, archived)
	if err != nil {
		return nil, fmt.Errorf("list units: %w", err)
	}
	return pgx.CollectRows(rows, scanUnit)
}

func (UnitStore) ByID(ctx context.Context, tx pgx.Tx, id string) (Unit, error) {
	rows, err := tx.Query(ctx, `SELECT `+unitColumns+` FROM core.units WHERE id::text = $1`, id)
	if err != nil {
		return Unit{}, fmt.Errorf("find unit: %w", err)
	}
	u, err := pgx.CollectExactlyOneRow(rows, scanUnit)
	if errors.Is(err, pgx.ErrNoRows) {
		return Unit{}, ErrUnitNotFound
	}
	return u, err
}

func (UnitStore) ByCode(ctx context.Context, tx pgx.Tx, code string) (Unit, error) {
	rows, err := tx.Query(ctx, `SELECT `+unitColumns+` FROM core.units WHERE code = $1`, code)
	if err != nil {
		return Unit{}, fmt.Errorf("find unit: %w", err)
	}
	u, err := pgx.CollectExactlyOneRow(rows, scanUnit)
	if errors.Is(err, pgx.ErrNoRows) {
		return Unit{}, ErrUnitNotFound
	}
	return u, err
}

func (UnitStore) Path(ctx context.Context, tx pgx.Tx, id string) ([]Unit, error) {
	rows, err := tx.Query(ctx, `
		WITH RECURSIVE up AS (
			SELECT u.*, 0 AS depth FROM core.units u WHERE u.id::text = $1
			UNION ALL
			SELECT p.*, up.depth + 1 FROM core.units p JOIN up ON p.id = up.parent_id
		)
		SELECT `+unitColumns+` FROM up ORDER BY depth DESC`, id)
	if err != nil {
		return nil, fmt.Errorf("unit path: %w", err)
	}
	return pgx.CollectRows(rows, scanUnit)
}

func (UnitStore) Insert(ctx context.Context, tx pgx.Tx, in UnitInput) (Unit, error) {
	rows, err := tx.Query(ctx, `
		INSERT INTO core.units (parent_id, kind, name, code) VALUES ($1, $2, $3, $4)
		RETURNING `+unitColumns, in.ParentID, in.Kind, in.Name, in.Code)
	if err != nil {
		return Unit{}, unitError(err)
	}
	u, err := pgx.CollectExactlyOneRow(rows, scanUnit)
	return u, unitError(err)
}

func (UnitStore) Update(ctx context.Context, tx pgx.Tx, id string, in UnitInput) (Unit, error) {
	rows, err := tx.Query(ctx, `
		UPDATE core.units SET kind = $2, name = $3, code = $4, updated_at = now()
		WHERE id::text = $1 RETURNING `+unitColumns, id, in.Kind, in.Name, in.Code)
	if err != nil {
		return Unit{}, unitError(err)
	}
	u, err := pgx.CollectExactlyOneRow(rows, scanUnit)
	return u, unitError(err)
}

func (UnitStore) SetParent(ctx context.Context, tx pgx.Tx, id string, parentID *string) error {
	_, err := tx.Exec(ctx, `UPDATE core.units SET parent_id = $2, updated_at = now() WHERE id::text = $1`, id, parentID)
	return unitError(err)
}

func (UnitStore) SetArchived(ctx context.Context, tx pgx.Tx, id string, archived bool) error {
	_, err := tx.Exec(ctx, `
		WITH RECURSIVE down AS (
			SELECT id FROM core.units WHERE id::text = $1
			UNION ALL
			SELECT c.id FROM core.units c JOIN down ON c.parent_id = down.id
		)
		UPDATE core.units SET archived_at = CASE WHEN $2 THEN coalesce(archived_at, now()) END, updated_at = now()
		WHERE id IN (SELECT id FROM down)`, id, archived)
	return err
}

func (UnitStore) IsInside(ctx context.Context, tx pgx.Tx, id, ancestorID string) (bool, error) {
	var inside bool
	err := tx.QueryRow(ctx, `
		WITH RECURSIVE up AS (
			SELECT id, parent_id FROM core.units WHERE id::text = $1
			UNION ALL
			SELECT p.id, p.parent_id FROM core.units p JOIN up ON p.id = up.parent_id
		)
		SELECT EXISTS (SELECT 1 FROM up WHERE id::text = $2)`, id, ancestorID).Scan(&inside)
	return inside, err
}

func scanUnit(row pgx.CollectableRow) (Unit, error) {
	var u Unit
	err := row.Scan(&u.ID, &u.ParentID, &u.Kind, &u.Name, &u.Code, &u.ArchivedAt, &u.CreatedAt)
	return u, err
}

func unitError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case uniqueViolation:
			return ErrCodeTaken
		case checkViolation:
			return ErrInvalidUnit
		case foreignKeyViolation, invalidText:
			return ErrUnitNotFound
		}
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrUnitNotFound
	}
	if err != nil {
		return fmt.Errorf("save unit: %w", err)
	}
	return nil
}
