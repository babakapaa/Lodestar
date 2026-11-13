package league

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

const awardColumns = `a.id::text, a.title, a.about, a.gpu_hours, a.quantity,
	(SELECT count(*)::int FROM league.grants g WHERE g.award_id = a.id), a.archived_at, a.created_at`

type AwardStore struct{}

func (AwardStore) Insert(ctx context.Context, tx pgx.Tx, authorID string, in AwardInput) (Award, error) {
	rows, err := tx.Query(ctx, `
		INSERT INTO league.awards AS a (title, about, gpu_hours, quantity, archived_at, created_by)
		VALUES ($1, $2, $3, $4, CASE WHEN $5 THEN now() END, $6) RETURNING `+awardColumns,
		in.Title, in.About, in.GPUHours, in.Quantity, in.Archived, authorID)
	if err != nil {
		return Award{}, fmt.Errorf("create award: %w", err)
	}
	return pgx.CollectExactlyOneRow(rows, scanAward)
}

func (AwardStore) Update(ctx context.Context, tx pgx.Tx, id string, in AwardInput) (Award, error) {
	rows, err := tx.Query(ctx, `
		UPDATE league.awards a SET title = $2, about = $3, gpu_hours = $4, quantity = $5,
		       archived_at = CASE WHEN $6 THEN coalesce(a.archived_at, now()) END
		WHERE a.id::text = $1 RETURNING `+awardColumns, id, in.Title, in.About, in.GPUHours, in.Quantity, in.Archived)
	if err != nil {
		return Award{}, fmt.Errorf("update award: %w", err)
	}
	a, err := pgx.CollectExactlyOneRow(rows, scanAward)
	if errors.Is(err, pgx.ErrNoRows) {
		return Award{}, ErrAwardNotFound
	}
	return a, err
}

func (AwardStore) ByID(ctx context.Context, tx pgx.Tx, id string) (Award, error) {
	rows, err := tx.Query(ctx, `SELECT `+awardColumns+` FROM league.awards a WHERE a.id::text = $1 FOR UPDATE`, id)
	if err != nil {
		return Award{}, fmt.Errorf("load award: %w", err)
	}
	a, err := pgx.CollectExactlyOneRow(rows, scanAward)
	if errors.Is(err, pgx.ErrNoRows) {
		return Award{}, ErrAwardNotFound
	}
	return a, err
}

func (AwardStore) All(ctx context.Context, tx pgx.Tx) ([]Award, error) {
	rows, err := tx.Query(ctx, `SELECT `+awardColumns+` FROM league.awards a ORDER BY a.archived_at NULLS FIRST, a.created_at, a.id`)
	if err != nil {
		return nil, fmt.Errorf("list awards: %w", err)
	}
	return pgx.CollectRows(rows, scanAward)
}

func scanAward(row pgx.CollectableRow) (Award, error) {
	var a Award
	err := row.Scan(&a.ID, &a.Title, &a.About, &a.GPUHours, &a.Quantity, &a.Granted, &a.ArchivedAt, &a.CreatedAt)
	return a, err
}
