package navigator

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type PassageStore struct{}

func (PassageStore) Pending(ctx context.Context, tx pgx.Tx, model string, limit int) ([]passage, error) {
	rows, err := tx.Query(ctx, `
		SELECT material_id::text, seq, text FROM navigator.passages
		WHERE model IS DISTINCT FROM $1 ORDER BY material_id, seq LIMIT $2`, model, limit)
	if err != nil {
		return nil, fmt.Errorf("pending passages: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (passage, error) {
		var p passage
		err := row.Scan(&p.MaterialID, &p.Seq, &p.Text)
		return p, err
	})
}

func (PassageStore) Embed(ctx context.Context, tx pgx.Tx, model string, p passage, v []float32) error {
	if _, err := tx.Exec(ctx, `
		UPDATE navigator.passages SET model = $1, embedding = $2::vector
		WHERE material_id::text = $3 AND seq = $4`, model, vectorLiteral(v), p.MaterialID, p.Seq); err != nil {
		return fmt.Errorf("save embedding: %w", err)
	}
	return nil
}

func (PassageStore) Search(ctx context.Context, tx pgx.Tx, model, courseID string, v []float32, limit int) ([]Fragment, error) {
	rows, err := tx.Query(ctx, `
		SELECT m.id::text, p.seq, m.title, m.reference, m.url, p.text
		FROM navigator.passages p JOIN navigator.materials m ON m.id = p.material_id
		WHERE m.course_id::text = $1 AND p.model = $2 AND vector_dims(p.embedding) = vector_dims($3::vector)
		ORDER BY p.embedding <=> $3::vector LIMIT $4`, courseID, model, vectorLiteral(v), limit)
	if err != nil {
		return nil, fmt.Errorf("search passages: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (Fragment, error) {
		var f Fragment
		err := row.Scan(&f.MaterialID, &f.Seq, &f.Title, &f.Reference, &f.URL, &f.Text)
		return f, err
	})
}
