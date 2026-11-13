package navigator

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type MaterialStore struct{}

const materialColumns = `m.id::text, m.course_id::text, m.title, m.reference, m.url,
	(SELECT count(*) FROM navigator.passages p WHERE p.material_id = m.id),
	(SELECT count(*) FROM navigator.passages p WHERE p.material_id = m.id AND p.model = $1),
	m.created_at`

func (s MaterialStore) Insert(ctx context.Context, tx pgx.Tx, model, courseID, authorID string, in MaterialInput, passages []string) (Material, error) {
	var id string
	err := tx.QueryRow(ctx, `
		INSERT INTO navigator.materials (course_id, author_id, title, reference, url) VALUES ($1, $2, $3, $4, $5)
		RETURNING id::text`, courseID, authorID, in.Title, in.Reference, in.URL).Scan(&id)
	if err != nil {
		return Material{}, fmt.Errorf("save material: %w", err)
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO navigator.passages (material_id, seq, text)
		SELECT $1, n, t FROM unnest($2::text[]) WITH ORDINALITY AS u(t, n)`, id, passages); err != nil {
		return Material{}, fmt.Errorf("save passages: %w", err)
	}
	return s.ByID(ctx, tx, model, id)
}

func (MaterialStore) ByID(ctx context.Context, tx pgx.Tx, model, id string) (Material, error) {
	rows, err := tx.Query(ctx, `SELECT `+materialColumns+` FROM navigator.materials m WHERE m.id::text = $2`, model, id)
	if err != nil {
		return Material{}, fmt.Errorf("load material: %w", err)
	}
	m, err := pgx.CollectExactlyOneRow(rows, scanMaterial)
	if errors.Is(err, pgx.ErrNoRows) {
		return Material{}, ErrMaterialNotFound
	}
	return m, err
}

func (MaterialStore) List(ctx context.Context, tx pgx.Tx, model, courseID string) ([]Material, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+materialColumns+` FROM navigator.materials m
		WHERE m.course_id::text = $2 ORDER BY m.created_at, m.id`, model, courseID)
	if err != nil {
		return nil, fmt.Errorf("list materials: %w", err)
	}
	return pgx.CollectRows(rows, scanMaterial)
}

func (MaterialStore) Delete(ctx context.Context, tx pgx.Tx, id string) error {
	if _, err := tx.Exec(ctx, `DELETE FROM navigator.materials WHERE id::text = $1`, id); err != nil {
		return fmt.Errorf("delete material: %w", err)
	}
	return nil
}

func scanMaterial(row pgx.CollectableRow) (Material, error) {
	var m Material
	err := row.Scan(&m.ID, &m.CourseID, &m.Title, &m.Reference, &m.URL, &m.Passages, &m.Indexed, &m.CreatedAt)
	return m, err
}
