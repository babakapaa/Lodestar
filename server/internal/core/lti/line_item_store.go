package lti

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type LineItemStore struct{}

func (LineItemStore) ByContext(ctx context.Context, tx pgx.Tx, contextID string) ([]LineItemLink, error) {
	rows, err := tx.Query(ctx, `
		SELECT id::text, url, work_kind, work_id FROM lti.line_items WHERE context_id::text = $1 ORDER BY created_at, id`, contextID)
	if err != nil {
		return nil, fmt.Errorf("lti line items: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (LineItemLink, error) {
		var l LineItemLink
		err := row.Scan(&l.ID, &l.URL, &l.WorkKind, &l.WorkID)
		return l, err
	})
}

func (LineItemStore) Put(ctx context.Context, tx pgx.Tx, contextID, url, kind, workID string) (LineItemLink, error) {
	l := LineItemLink{URL: url, WorkKind: kind, WorkID: workID}
	err := tx.QueryRow(ctx, `
		INSERT INTO lti.line_items (context_id, url, work_kind, work_id) VALUES ($1, $2, $3, $4)
		ON CONFLICT (context_id, url) DO UPDATE SET work_kind = EXCLUDED.work_kind, work_id = EXCLUDED.work_id
		RETURNING id::text`, contextID, url, kind, workID).Scan(&l.ID)
	if err != nil {
		return LineItemLink{}, fmt.Errorf("save lti line item: %w", err)
	}
	return l, nil
}
