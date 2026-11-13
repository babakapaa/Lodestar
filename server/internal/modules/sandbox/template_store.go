package sandbox

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const templateColumns = `t.id::text, t.name, t.about, t.image, t.archived_at IS NOT NULL, t.created_at`

type TemplateStore struct{}

func (TemplateStore) Insert(ctx context.Context, tx pgx.Tx, in TemplateInput) (Template, error) {
	rows, err := tx.Query(ctx, `
		INSERT INTO sandbox.templates AS t (name, about, image, archived_at)
		VALUES ($1, $2, $3, CASE WHEN $4 THEN now() END) RETURNING `+templateColumns, in.Name, in.About, in.Image, in.Archived)
	if err != nil {
		return Template{}, fmt.Errorf("create template: %w", err)
	}
	t, err := pgx.CollectExactlyOneRow(rows, scanTemplate)
	return t, templateTaken(err)
}

func (TemplateStore) Update(ctx context.Context, tx pgx.Tx, id string, in TemplateInput) (Template, error) {
	rows, err := tx.Query(ctx, `
		UPDATE sandbox.templates t SET name = $2, about = $3, image = $4,
		       archived_at = CASE WHEN $5 THEN coalesce(t.archived_at, now()) END
		WHERE t.id::text = $1 RETURNING `+templateColumns, id, in.Name, in.About, in.Image, in.Archived)
	if err != nil {
		return Template{}, fmt.Errorf("update template: %w", err)
	}
	t, err := pgx.CollectExactlyOneRow(rows, scanTemplate)
	if errors.Is(err, pgx.ErrNoRows) {
		return Template{}, ErrTemplateNotFound
	}
	return t, templateTaken(err)
}

func (TemplateStore) ByID(ctx context.Context, tx pgx.Tx, id string) (Template, error) {
	rows, err := tx.Query(ctx, `SELECT `+templateColumns+` FROM sandbox.templates t WHERE t.id::text = $1`, id)
	if err != nil {
		return Template{}, fmt.Errorf("load template: %w", err)
	}
	t, err := pgx.CollectExactlyOneRow(rows, scanTemplate)
	if errors.Is(err, pgx.ErrNoRows) {
		return Template{}, ErrTemplateNotFound
	}
	return t, err
}

func (TemplateStore) All(ctx context.Context, tx pgx.Tx) ([]Template, error) {
	rows, err := tx.Query(ctx, `SELECT `+templateColumns+` FROM sandbox.templates t ORDER BY t.archived_at NULLS FIRST, t.name`)
	if err != nil {
		return nil, fmt.Errorf("list templates: %w", err)
	}
	return pgx.CollectRows(rows, scanTemplate)
}

func scanTemplate(row pgx.CollectableRow) (Template, error) {
	var t Template
	err := row.Scan(&t.ID, &t.Name, &t.About, &t.Image, &t.Archived, &t.CreatedAt)
	return t, err
}

func templateTaken(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrTemplateTaken
	}
	return err
}
