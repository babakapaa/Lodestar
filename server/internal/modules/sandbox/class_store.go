package sandbox

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const classColumns = `c.id::text, c.name, c.memory_gb, c.rate::float8, c.archived_at IS NOT NULL, c.created_at`

type ClassStore struct{}

func (ClassStore) Insert(ctx context.Context, tx pgx.Tx, in GPUClassInput) (GPUClass, error) {
	rows, err := tx.Query(ctx, `
		INSERT INTO sandbox.gpu_classes AS c (name, memory_gb, rate, archived_at)
		VALUES ($1, $2, $3, CASE WHEN $4 THEN now() END) RETURNING `+classColumns,
		in.Name, in.MemoryGB, in.Rate, in.Archived)
	if err != nil {
		return GPUClass{}, fmt.Errorf("create gpu class: %w", err)
	}
	c, err := pgx.CollectExactlyOneRow(rows, scanClass)
	return c, taken(err)
}

func (ClassStore) Update(ctx context.Context, tx pgx.Tx, id string, in GPUClassInput) (GPUClass, error) {
	rows, err := tx.Query(ctx, `
		UPDATE sandbox.gpu_classes c SET name = $2, memory_gb = $3, rate = $4,
		       archived_at = CASE WHEN $5 THEN coalesce(c.archived_at, now()) END
		WHERE c.id::text = $1 RETURNING `+classColumns, id, in.Name, in.MemoryGB, in.Rate, in.Archived)
	if err != nil {
		return GPUClass{}, fmt.Errorf("update gpu class: %w", err)
	}
	c, err := pgx.CollectExactlyOneRow(rows, scanClass)
	if errors.Is(err, pgx.ErrNoRows) {
		return GPUClass{}, ErrClassNotFound
	}
	return c, taken(err)
}

func (ClassStore) All(ctx context.Context, tx pgx.Tx) ([]GPUClass, error) {
	rows, err := tx.Query(ctx, `SELECT `+classColumns+` FROM sandbox.gpu_classes c ORDER BY c.archived_at NULLS FIRST, c.rate, c.name`)
	if err != nil {
		return nil, fmt.Errorf("list gpu classes: %w", err)
	}
	return pgx.CollectRows(rows, scanClass)
}

func scanClass(row pgx.CollectableRow) (GPUClass, error) {
	var c GPUClass
	err := row.Scan(&c.ID, &c.Name, &c.MemoryGB, &c.Rate, &c.Archived, &c.CreatedAt)
	return c, err
}

func taken(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrClassTaken
	}
	return err
}
