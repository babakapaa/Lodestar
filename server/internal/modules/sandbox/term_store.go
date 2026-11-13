package sandbox

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

const termColumns = `t.id::text, t.name, t.starts_on::text, t.ends_on::text`

type TermStore struct{}

func (TermStore) Lock(ctx context.Context, tx pgx.Tx) error {
	if _, err := tx.Exec(ctx, `SELECT pg_advisory_xact_lock(hashtext('sandbox.terms:' || current_setting('lodestar.tenant')))`); err != nil {
		return fmt.Errorf("lock terms: %w", err)
	}
	return nil
}

func (TermStore) Overlaps(ctx context.Context, tx pgx.Tx, exceptID string, in TermInput) (bool, error) {
	var found bool
	err := tx.QueryRow(ctx, `
		SELECT EXISTS (SELECT 1 FROM sandbox.terms t
		WHERE t.id::text <> $1 AND t.starts_on <= $3::date AND t.ends_on >= $2::date)`,
		exceptID, in.StartsOn, in.EndsOn).Scan(&found)
	if err != nil {
		return false, fmt.Errorf("check term overlap: %w", err)
	}
	return found, nil
}

func (TermStore) Insert(ctx context.Context, tx pgx.Tx, in TermInput) (Term, error) {
	rows, err := tx.Query(ctx, `
		INSERT INTO sandbox.terms AS t (name, starts_on, ends_on) VALUES ($1, $2::date, $3::date)
		RETURNING `+termColumns, in.Name, in.StartsOn, in.EndsOn)
	if err != nil {
		return Term{}, fmt.Errorf("create term: %w", err)
	}
	return pgx.CollectExactlyOneRow(rows, scanTerm)
}

func (TermStore) Update(ctx context.Context, tx pgx.Tx, id string, in TermInput) (Term, error) {
	rows, err := tx.Query(ctx, `
		UPDATE sandbox.terms t SET name = $2, starts_on = $3::date, ends_on = $4::date
		WHERE t.id::text = $1 RETURNING `+termColumns, id, in.Name, in.StartsOn, in.EndsOn)
	if err != nil {
		return Term{}, fmt.Errorf("update term: %w", err)
	}
	t, err := pgx.CollectExactlyOneRow(rows, scanTerm)
	if errors.Is(err, pgx.ErrNoRows) {
		return Term{}, ErrTermNotFound
	}
	return t, err
}

func (TermStore) Delete(ctx context.Context, tx pgx.Tx, id string) error {
	tag, err := tx.Exec(ctx, `DELETE FROM sandbox.terms WHERE id::text = $1`, id)
	if err != nil {
		return fmt.Errorf("delete term: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrTermNotFound
	}
	return nil
}

func (TermStore) All(ctx context.Context, tx pgx.Tx) ([]Term, error) {
	rows, err := tx.Query(ctx, `SELECT `+termColumns+` FROM sandbox.terms t ORDER BY t.starts_on DESC`)
	if err != nil {
		return nil, fmt.Errorf("list terms: %w", err)
	}
	return pgx.CollectRows(rows, scanTerm)
}

func (TermStore) Current(ctx context.Context, tx pgx.Tx) (Term, bool, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+termColumns+` FROM sandbox.terms t
		WHERE current_date BETWEEN t.starts_on AND t.ends_on ORDER BY t.starts_on DESC LIMIT 1`)
	if err != nil {
		return Term{}, false, fmt.Errorf("current term: %w", err)
	}
	t, err := pgx.CollectExactlyOneRow(rows, scanTerm)
	if errors.Is(err, pgx.ErrNoRows) {
		return Term{}, false, nil
	}
	return t, err == nil, err
}

func scanTerm(row pgx.CollectableRow) (Term, error) {
	var t Term
	err := row.Scan(&t.ID, &t.Name, &t.StartsOn, &t.EndsOn)
	return t, err
}
