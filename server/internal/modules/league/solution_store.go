package league

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

const solutionColumns = `s.id::text, s.task_id::text, s.team_id, s.team_name, s.proposed_by::text, s.members::text[], s.work, s.votes,
	s.similar_to::text, s.similarity, s.explanation, s.points, s.passed::text[], s.comment, s.reviewed_at, s.created_at`

type SolutionStore struct{}

func (SolutionStore) Insert(ctx context.Context, tx pgx.Tx, proposalID string, s Solution) (Solution, bool, error) {
	rows, err := tx.Query(ctx, `
		INSERT INTO league.solutions AS s (task_id, team_id, team_name, proposal_id, proposed_by, members, work, votes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (task_id, proposal_id) DO NOTHING RETURNING `+solutionColumns,
		s.TaskID, s.TeamID, s.TeamName, proposalID, s.ProposedBy, s.Members, s.Work, s.Votes)
	if err != nil {
		return Solution{}, false, fmt.Errorf("save solution: %w", err)
	}
	saved, err := pgx.CollectExactlyOneRow(rows, scanSolution)
	if errors.Is(err, pgx.ErrNoRows) {
		return Solution{}, false, nil
	}
	if err != nil {
		return Solution{}, false, fmt.Errorf("save solution: %w", err)
	}
	return saved, true, nil
}

func (SolutionStore) ByID(ctx context.Context, tx pgx.Tx, id string) (Solution, error) {
	rows, err := tx.Query(ctx, `SELECT `+solutionColumns+` FROM league.solutions s WHERE s.id::text = $1 FOR UPDATE`, id)
	if err != nil {
		return Solution{}, fmt.Errorf("load solution: %w", err)
	}
	s, err := pgx.CollectExactlyOneRow(rows, scanSolution)
	if errors.Is(err, pgx.ErrNoRows) {
		return Solution{}, ErrSolutionNotFound
	}
	return s, err
}

func (SolutionStore) ByTask(ctx context.Context, tx pgx.Tx, taskID string) ([]Solution, error) {
	rows, err := tx.Query(ctx, `SELECT `+solutionColumns+` FROM league.solutions s WHERE s.task_id::text = $1 ORDER BY s.created_at, s.id`, taskID)
	if err != nil {
		return nil, fmt.Errorf("list solutions: %w", err)
	}
	return pgx.CollectRows(rows, scanSolution)
}

func (SolutionStore) Embeddings(ctx context.Context, tx pgx.Tx, exceptTeam string) (map[string][]float32, error) {
	rows, err := tx.Query(ctx, `SELECT id::text, embedding FROM league.solutions WHERE embedding IS NOT NULL AND team_id <> $1`, exceptTeam)
	if err != nil {
		return nil, fmt.Errorf("load embeddings: %w", err)
	}
	out := map[string][]float32{}
	var id string
	var vector []float32
	_, err = pgx.ForEachRow(rows, []any{&id, &vector}, func() error {
		out[id] = vector
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("load embeddings: %w", err)
	}
	return out, nil
}

func (SolutionStore) Embed(ctx context.Context, tx pgx.Tx, id string, vector []float32, similarTo *string, similarity *float64) error {
	_, err := tx.Exec(ctx, `UPDATE league.solutions SET embedding = $2, similar_to = $3, similarity = $4 WHERE id::text = $1`, id, vector, similarTo, similarity)
	if err != nil {
		return fmt.Errorf("save embedding: %w", err)
	}
	return nil
}

func (SolutionStore) Similar(ctx context.Context, tx pgx.Tx, id string) (SimilarSolution, error) {
	var out SimilarSolution
	var work SolutionWork
	err := tx.QueryRow(ctx, `
		SELECT s.id::text, t.title, s.team_name, s.work FROM league.solutions s JOIN league.tasks t ON t.id = s.task_id
		WHERE s.id::text = $1`, id).Scan(&out.ID, &out.TaskTitle, &out.TeamName, &work)
	if err != nil {
		return out, fmt.Errorf("load similar solution: %w", err)
	}
	out.Text = work.Text
	return out, nil
}

func (SolutionStore) Review(ctx context.Context, tx pgx.Tx, id, reviewerID string, in ReviewInput) (Solution, error) {
	rows, err := tx.Query(ctx, `
		UPDATE league.solutions s SET points = $2, passed = $3, comment = $4, reviewed_by = $5, reviewed_at = now()
		WHERE s.id::text = $1 RETURNING `+solutionColumns, id, in.Points, in.Passed, in.Comment, reviewerID)
	if err != nil {
		return Solution{}, fmt.Errorf("review solution: %w", err)
	}
	return pgx.CollectExactlyOneRow(rows, scanSolution)
}

func (SolutionStore) Explain(ctx context.Context, tx pgx.Tx, id, text string) (Solution, error) {
	rows, err := tx.Query(ctx, `UPDATE league.solutions s SET explanation = $2 WHERE s.id::text = $1 RETURNING `+solutionColumns, id, text)
	if err != nil {
		return Solution{}, fmt.Errorf("explain solution: %w", err)
	}
	return pgx.CollectExactlyOneRow(rows, scanSolution)
}

func scanSolution(row pgx.CollectableRow) (Solution, error) {
	var s Solution
	err := row.Scan(&s.ID, &s.TaskID, &s.TeamID, &s.TeamName, &s.ProposedBy, &s.Members, &s.Work, &s.Votes,
		&s.SimilarTo, &s.Similarity, &s.Explanation, &s.Points, &s.Passed, &s.Comment, &s.ReviewedAt, &s.CreatedAt)
	return s, err
}
