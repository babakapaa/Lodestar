package league

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

const entryColumns = `e.id::text, e.season_id::text, e.user_id::text, e.kind, e.work_kind, e.work_id, e.title,
	e.score, e.max, e.weight, e.points, e.reason, e.created_at, e.updated_at`

type EntryStore struct{}

func (EntryStore) Works(ctx context.Context, tx pgx.Tx, leagueID string) (map[string]Entry, error) {
	rows, err := tx.Query(ctx, `SELECT `+entryColumns+` FROM league.entries e WHERE e.league_id::text = $1 AND e.kind = 'work'`, leagueID)
	if err != nil {
		return nil, fmt.Errorf("list work entries: %w", err)
	}
	list, err := pgx.CollectRows(rows, scanEntry)
	if err != nil {
		return nil, fmt.Errorf("list work entries: %w", err)
	}
	out := make(map[string]Entry, len(list))
	for _, e := range list {
		out[workKey(e.UserID, *e.WorkKind, *e.WorkID)] = e
	}
	return out, nil
}

func (EntryStore) InsertWork(ctx context.Context, tx pgx.Tx, leagueID, seasonID string, e Entry) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO league.entries (league_id, season_id, user_id, kind, work_kind, work_id, title, score, max, weight, points)
		VALUES ($1, $2, $3, 'work', $4, $5, $6, $7, $8, $9, $10)`,
		leagueID, seasonID, e.UserID, e.WorkKind, e.WorkID, e.Title, e.Score, e.Max, e.Weight, e.Points)
	if err != nil {
		return fmt.Errorf("add work entry: %w", err)
	}
	return nil
}

func (EntryStore) UpdateWork(ctx context.Context, tx pgx.Tx, e Entry) error {
	_, err := tx.Exec(ctx, `
		UPDATE league.entries SET title = $2, score = $3, max = $4, weight = $5, points = $6, updated_at = now()
		WHERE id::text = $1`, e.ID, e.Title, e.Score, e.Max, e.Weight, e.Points)
	if err != nil {
		return fmt.Errorf("update work entry: %w", err)
	}
	return nil
}

func (EntryStore) Delete(ctx context.Context, tx pgx.Tx, id string) error {
	if _, err := tx.Exec(ctx, `DELETE FROM league.entries WHERE id::text = $1`, id); err != nil {
		return fmt.Errorf("delete entry: %w", err)
	}
	return nil
}

func (EntryStore) Adjust(ctx context.Context, tx pgx.Tx, leagueID, seasonID, userID string, points int, reason, authorID string) (Entry, error) {
	rows, err := tx.Query(ctx, `
		INSERT INTO league.entries AS e (league_id, season_id, user_id, kind, points, reason, created_by)
		VALUES ($1, $2, $3, 'adjustment', $4, $5, $6) RETURNING `+entryColumns,
		leagueID, seasonID, userID, points, reason, authorID)
	if err != nil {
		return Entry{}, fmt.Errorf("adjust points: %w", err)
	}
	return pgx.CollectExactlyOneRow(rows, scanEntry)
}

func (EntryStore) Totals(ctx context.Context, tx pgx.Tx, seasonID string) (map[string]int, error) {
	rows, err := tx.Query(ctx, `SELECT user_id::text, sum(points)::int FROM league.entries WHERE season_id::text = $1 GROUP BY user_id`, seasonID)
	if err != nil {
		return nil, fmt.Errorf("sum points: %w", err)
	}
	out := map[string]int{}
	var userID string
	var points int
	_, err = pgx.ForEachRow(rows, []any{&userID, &points}, func() error {
		out[userID] = points
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("sum points: %w", err)
	}
	return out, nil
}

func (EntryStore) Of(ctx context.Context, tx pgx.Tx, seasonID, userID string) ([]Entry, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+entryColumns+` FROM league.entries e
		WHERE e.season_id::text = $1 AND e.user_id::text = $2 ORDER BY e.created_at, e.id`, seasonID, userID)
	if err != nil {
		return nil, fmt.Errorf("list entries: %w", err)
	}
	return pgx.CollectRows(rows, scanEntry)
}

func (EntryStore) ByID(ctx context.Context, tx pgx.Tx, leagueID, id string) (Entry, error) {
	rows, err := tx.Query(ctx, `SELECT `+entryColumns+` FROM league.entries e WHERE e.league_id::text = $1 AND e.id::text = $2`, leagueID, id)
	if err != nil {
		return Entry{}, fmt.Errorf("load entry: %w", err)
	}
	e, err := pgx.CollectExactlyOneRow(rows, scanEntry)
	if errors.Is(err, pgx.ErrNoRows) {
		return Entry{}, ErrDisputeTarget
	}
	return e, err
}

func scanEntry(row pgx.CollectableRow) (Entry, error) {
	var e Entry
	err := row.Scan(&e.ID, &e.SeasonID, &e.UserID, &e.Kind, &e.WorkKind, &e.WorkID, &e.Title,
		&e.Score, &e.Max, &e.Weight, &e.Points, &e.Reason, &e.CreatedAt, &e.UpdatedAt)
	return e, err
}

func workKey(userID, kind, id string) string {
	return userID + "/" + kind + "/" + id
}
