package study

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type SearchMapPrivacySection struct{}

func (SearchMapPrivacySection) Name() string {
	return "study.search_maps"
}

func (SearchMapPrivacySection) Export(ctx context.Context, tx pgx.Tx, userID string) (any, error) {
	rows, err := tx.Query(ctx, `
		SELECT a.title, m.data, m.changed_at FROM study.search_maps m
		JOIN study.assignments a ON a.id = m.assignment_id
		WHERE m.student_id::text = $1 ORDER BY m.changed_at`, userID)
	if err != nil {
		return nil, fmt.Errorf("export search maps: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (map[string]any, error) {
		var title string
		var m SearchMap
		var at time.Time
		err := row.Scan(&title, &m, &at)
		return map[string]any{"assignment": title, "map": m, "changedAt": at}, err
	})
}
