package league

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

const seasonColumns = `s.id::text, s.number, s.started_at, s.ended_at`

type SeasonStore struct{}

func (SeasonStore) Current(ctx context.Context, tx pgx.Tx, leagueID string) (Season, error) {
	_, err := tx.Exec(ctx, `
		INSERT INTO league.seasons (league_id, number)
		SELECT $1, coalesce(max(number), 0) + 1 FROM league.seasons WHERE league_id = $1
		HAVING count(*) FILTER (WHERE ended_at IS NULL) = 0`, leagueID)
	if err != nil {
		return Season{}, fmt.Errorf("start season: %w", err)
	}
	rows, err := tx.Query(ctx, `SELECT `+seasonColumns+` FROM league.seasons s WHERE s.league_id::text = $1 AND s.ended_at IS NULL`, leagueID)
	if err != nil {
		return Season{}, fmt.Errorf("load season: %w", err)
	}
	return pgx.CollectExactlyOneRow(rows, scanSeason)
}

func (SeasonStore) Number(ctx context.Context, tx pgx.Tx, leagueID string, number int) (*Season, error) {
	rows, err := tx.Query(ctx, `SELECT `+seasonColumns+` FROM league.seasons s WHERE s.league_id::text = $1 AND s.number = $2`, leagueID, number)
	if err != nil {
		return nil, fmt.Errorf("load season: %w", err)
	}
	s, err := pgx.CollectExactlyOneRow(rows, scanSeason)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("load season: %w", err)
	}
	return &s, nil
}

func (SeasonStore) End(ctx context.Context, tx pgx.Tx, id string) error {
	if _, err := tx.Exec(ctx, `UPDATE league.seasons SET ended_at = now() WHERE id::text = $1`, id); err != nil {
		return fmt.Errorf("end season: %w", err)
	}
	return nil
}

func scanSeason(row pgx.CollectableRow) (Season, error) {
	var s Season
	err := row.Scan(&s.ID, &s.Number, &s.StartedAt, &s.EndedAt)
	return s, err
}
