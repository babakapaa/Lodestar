package study

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type SearchMapStore struct{}

func (SearchMapStore) Get(ctx context.Context, tx pgx.Tx, assignmentID, studentID string) (SearchMap, *time.Time, error) {
	var m SearchMap
	var at time.Time
	err := tx.QueryRow(ctx, `
		SELECT data, changed_at FROM study.search_maps
		WHERE assignment_id::text = $1 AND student_id::text = $2`, assignmentID, studentID).Scan(&m, &at)
	if errors.Is(err, pgx.ErrNoRows) {
		return SearchMap{Used: []MapSource{}, Rejected: []MapSource{}, Queries: []string{}}, nil, nil
	}
	if err != nil {
		return SearchMap{}, nil, fmt.Errorf("load search map: %w", err)
	}
	return m, &at, nil
}

func (SearchMapStore) Put(ctx context.Context, tx pgx.Tx, assignmentID, studentID string, m SearchMap) (time.Time, error) {
	var at time.Time
	err := tx.QueryRow(ctx, `
		INSERT INTO study.search_maps (assignment_id, student_id, data) VALUES ($1, $2, $3)
		ON CONFLICT (assignment_id, student_id) DO UPDATE SET data = excluded.data, changed_at = now()
		RETURNING changed_at`, assignmentID, studentID, m).Scan(&at)
	if err != nil {
		return time.Time{}, fmt.Errorf("save search map: %w", err)
	}
	return at, nil
}
