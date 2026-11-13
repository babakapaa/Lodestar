package league

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

const taskColumns = `t.id::text, t.league_id::text, t.title, t.about, t.partner, t.points, t.min_size, t.max_size, t.self_formed, t.due_at, t.closed_at, t.created_at`

type TaskStore struct{}

func (TaskStore) Insert(ctx context.Context, tx pgx.Tx, leagueID, authorID string, in TaskInput) (Task, error) {
	rows, err := tx.Query(ctx, `
		INSERT INTO league.tasks AS t (league_id, title, about, partner, points, min_size, max_size, self_formed, due_at, closed_at, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, CASE WHEN $10 THEN now() END, $11) RETURNING `+taskColumns,
		leagueID, in.Title, in.About, in.Partner, in.Points, in.MinSize, in.MaxSize, in.SelfFormed, in.DueAt, in.Closed, authorID)
	if err != nil {
		return Task{}, fmt.Errorf("create task: %w", err)
	}
	return pgx.CollectExactlyOneRow(rows, scanTask)
}

func (TaskStore) Update(ctx context.Context, tx pgx.Tx, id string, in TaskInput) (Task, error) {
	rows, err := tx.Query(ctx, `
		UPDATE league.tasks t SET title = $2, about = $3, partner = $4, points = $5, min_size = $6, max_size = $7,
		       self_formed = $8, due_at = $9, closed_at = CASE WHEN $10 THEN coalesce(t.closed_at, now()) END
		WHERE t.id::text = $1 RETURNING `+taskColumns,
		id, in.Title, in.About, in.Partner, in.Points, in.MinSize, in.MaxSize, in.SelfFormed, in.DueAt, in.Closed)
	if err != nil {
		return Task{}, fmt.Errorf("update task: %w", err)
	}
	return pgx.CollectExactlyOneRow(rows, scanTask)
}

func (TaskStore) ByID(ctx context.Context, tx pgx.Tx, id string) (Task, error) {
	rows, err := tx.Query(ctx, `SELECT `+taskColumns+` FROM league.tasks t WHERE t.id::text = $1`, id)
	if err != nil {
		return Task{}, fmt.Errorf("load task: %w", err)
	}
	t, err := pgx.CollectExactlyOneRow(rows, scanTask)
	if errors.Is(err, pgx.ErrNoRows) {
		return Task{}, ErrTaskNotFound
	}
	return t, err
}

func (TaskStore) ByLeague(ctx context.Context, tx pgx.Tx, leagueID string) ([]Task, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+taskColumns+` FROM league.tasks t WHERE t.league_id::text = $1
		ORDER BY t.closed_at NULLS FIRST, t.due_at NULLS LAST, t.created_at, t.id`, leagueID)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	return pgx.CollectRows(rows, scanTask)
}

func scanTask(row pgx.CollectableRow) (Task, error) {
	var t Task
	err := row.Scan(&t.ID, &t.LeagueID, &t.Title, &t.About, &t.Partner, &t.Points, &t.MinSize, &t.MaxSize, &t.SelfFormed, &t.DueAt, &t.ClosedAt, &t.CreatedAt)
	return t, err
}
