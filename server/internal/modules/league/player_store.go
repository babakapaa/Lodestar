package league

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type PlayerStore struct{}

func (PlayerStore) Ensure(ctx context.Context, tx pgx.Tx, leagueID string, userIDs []string) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO league.players (league_id, user_id) SELECT $1, u FROM unnest($2::uuid[]) u
		ON CONFLICT DO NOTHING`, leagueID, userIDs)
	if err != nil {
		return fmt.Errorf("add players: %w", err)
	}
	return nil
}

func (PlayerStore) All(ctx context.Context, tx pgx.Tx, leagueID string) (map[string]Player, error) {
	rows, err := tx.Query(ctx, `SELECT user_id::text, division, display, nickname, joined_at FROM league.players WHERE league_id::text = $1`, leagueID)
	if err != nil {
		return nil, fmt.Errorf("list players: %w", err)
	}
	list, err := pgx.CollectRows(rows, scanPlayer)
	if err != nil {
		return nil, fmt.Errorf("list players: %w", err)
	}
	out := make(map[string]Player, len(list))
	for _, p := range list {
		out[p.UserID] = p
	}
	return out, nil
}

func (PlayerStore) One(ctx context.Context, tx pgx.Tx, leagueID, userID string) (Player, error) {
	rows, err := tx.Query(ctx, `
		SELECT user_id::text, division, display, nickname, joined_at FROM league.players
		WHERE league_id::text = $1 AND user_id::text = $2`, leagueID, userID)
	if err != nil {
		return Player{}, fmt.Errorf("load player: %w", err)
	}
	p, err := pgx.CollectExactlyOneRow(rows, scanPlayer)
	if errors.Is(err, pgx.ErrNoRows) {
		return Player{}, ErrPlayerNotFound
	}
	return p, err
}

func (PlayerStore) SetDivision(ctx context.Context, tx pgx.Tx, leagueID, userID string, division int) error {
	_, err := tx.Exec(ctx, `UPDATE league.players SET division = $3 WHERE league_id::text = $1 AND user_id::text = $2`, leagueID, userID, division)
	if err != nil {
		return fmt.Errorf("move player: %w", err)
	}
	return nil
}

func (PlayerStore) SetDisplay(ctx context.Context, tx pgx.Tx, leagueID, userID string, in DisplayInput) error {
	_, err := tx.Exec(ctx, `UPDATE league.players SET display = $3, nickname = $4 WHERE league_id::text = $1 AND user_id::text = $2`,
		leagueID, userID, in.Display, in.Nickname)
	if err != nil {
		return fmt.Errorf("save display: %w", err)
	}
	return nil
}

func scanPlayer(row pgx.CollectableRow) (Player, error) {
	var p Player
	err := row.Scan(&p.UserID, &p.Division, &p.Display, &p.Nickname, &p.JoinedAt)
	return p, err
}

func (PlayerStore) Clamp(ctx context.Context, tx pgx.Tx, leagueID string, top int) error {
	_, err := tx.Exec(ctx, `UPDATE league.players SET division = $2 WHERE league_id::text = $1 AND division > $2`, leagueID, top)
	if err != nil {
		return fmt.Errorf("clamp divisions: %w", err)
	}
	return nil
}
