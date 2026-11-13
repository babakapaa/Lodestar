package league

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

const transitionColumns = `t.id::text, t.season_id::text, t.user_id::text, t.from_division, t.to_division, t.points, t.status, t.decided_at, t.created_at`

type TransitionStore struct{}

func (TransitionStore) Apply(ctx context.Context, tx pgx.Tx, leagueID string, t Transition) (Transition, error) {
	rows, err := tx.Query(ctx, `
		INSERT INTO league.transitions AS t (league_id, season_id, user_id, from_division, to_division, points, status, decided_at)
		VALUES ($1, $2, $3, $4, $5, $6, 'applied', now()) RETURNING `+transitionColumns,
		leagueID, t.SeasonID, t.UserID, t.From, t.To, t.Points)
	if err != nil {
		return Transition{}, fmt.Errorf("move player: %w", err)
	}
	return pgx.CollectExactlyOneRow(rows, scanTransition)
}

func (TransitionStore) Queue(ctx context.Context, tx pgx.Tx, leagueID string, t Transition) (Transition, bool, error) {
	var queued Transition
	var created bool
	err := tx.QueryRow(ctx, `
		INSERT INTO league.transitions AS t (league_id, season_id, user_id, from_division, to_division, points, status)
		VALUES ($1, $2, $3, $4, $5, $6, 'pending')
		ON CONFLICT (league_id, user_id) WHERE status = 'pending' DO UPDATE
		SET season_id = EXCLUDED.season_id, from_division = EXCLUDED.from_division,
		    to_division = EXCLUDED.to_division, points = EXCLUDED.points
		RETURNING `+transitionColumns+`, xmax = 0`,
		leagueID, t.SeasonID, t.UserID, t.From, t.To, t.Points).Scan(
		&queued.ID, &queued.SeasonID, &queued.UserID, &queued.From, &queued.To, &queued.Points,
		&queued.Status, &queued.DecidedAt, &queued.CreatedAt, &created)
	if err != nil {
		return Transition{}, false, fmt.Errorf("queue move: %w", err)
	}
	return queued, created, nil
}

func (TransitionStore) Highest(ctx context.Context, tx pgx.Tx, seasonID, userID string) (int, error) {
	var highest int
	err := tx.QueryRow(ctx, `
		SELECT coalesce(max(to_division), -1) FROM league.transitions
		WHERE season_id::text = $1 AND user_id::text = $2 AND to_division > from_division`, seasonID, userID).Scan(&highest)
	if err != nil {
		return 0, fmt.Errorf("highest move: %w", err)
	}
	return highest, nil
}

func (TransitionStore) ByID(ctx context.Context, tx pgx.Tx, id string) (Transition, string, error) {
	var t Transition
	var leagueID string
	err := tx.QueryRow(ctx, `SELECT `+transitionColumns+`, t.league_id::text FROM league.transitions t WHERE t.id::text = $1 FOR UPDATE`, id).Scan(
		&t.ID, &t.SeasonID, &t.UserID, &t.From, &t.To, &t.Points, &t.Status, &t.DecidedAt, &t.CreatedAt, &leagueID)
	if errors.Is(err, pgx.ErrNoRows) {
		return t, "", ErrTransitionNotFound
	}
	if err != nil {
		return t, "", fmt.Errorf("load move: %w", err)
	}
	return t, leagueID, nil
}

func (TransitionStore) Decide(ctx context.Context, tx pgx.Tx, id string, status TransitionStatus, deciderID string) (Transition, error) {
	rows, err := tx.Query(ctx, `
		UPDATE league.transitions t SET status = $2, decided_by = $3, decided_at = now()
		WHERE t.id::text = $1 RETURNING `+transitionColumns, id, status, deciderID)
	if err != nil {
		return Transition{}, fmt.Errorf("decide move: %w", err)
	}
	return pgx.CollectExactlyOneRow(rows, scanTransition)
}

func (TransitionStore) Pending(ctx context.Context, tx pgx.Tx, leagueID string) ([]Transition, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+transitionColumns+` FROM league.transitions t
		WHERE t.league_id::text = $1 AND t.status = 'pending' ORDER BY t.created_at, t.id`, leagueID)
	if err != nil {
		return nil, fmt.Errorf("pending moves: %w", err)
	}
	return pgx.CollectRows(rows, scanTransition)
}

func (TransitionStore) Of(ctx context.Context, tx pgx.Tx, leagueID, userID string) ([]Transition, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+transitionColumns+` FROM league.transitions t
		WHERE t.league_id::text = $1 AND t.user_id::text = $2 ORDER BY t.created_at, t.id`, leagueID, userID)
	if err != nil {
		return nil, fmt.Errorf("list moves: %w", err)
	}
	return pgx.CollectRows(rows, scanTransition)
}

func scanTransition(row pgx.CollectableRow) (Transition, error) {
	var t Transition
	err := row.Scan(&t.ID, &t.SeasonID, &t.UserID, &t.From, &t.To, &t.Points, &t.Status, &t.DecidedAt, &t.CreatedAt)
	return t, err
}
