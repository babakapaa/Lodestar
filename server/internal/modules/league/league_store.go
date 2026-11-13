package league

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

const leagueColumns = `l.id::text, l.origin_kind, l.origin_id, l.divisions, l.confirm, l.updated_at`

type LeagueStore struct{}

func (LeagueStore) Ensure(ctx context.Context, tx pgx.Tx, kind, id string) (League, error) {
	_, err := tx.Exec(ctx, `
		INSERT INTO league.leagues (origin_kind, origin_id, divisions) VALUES ($1, $2, $3)
		ON CONFLICT DO NOTHING`, kind, id, defaultDivisions())
	if err != nil {
		return League{}, fmt.Errorf("create league: %w", err)
	}
	rows, err := tx.Query(ctx, `SELECT `+leagueColumns+` FROM league.leagues l WHERE l.origin_kind = $1 AND l.origin_id = $2 FOR UPDATE`, kind, id)
	if err != nil {
		return League{}, fmt.Errorf("load league: %w", err)
	}
	return pgx.CollectExactlyOneRow(rows, scanLeague)
}

func (LeagueStore) ByID(ctx context.Context, tx pgx.Tx, id string) (League, error) {
	rows, err := tx.Query(ctx, `SELECT `+leagueColumns+` FROM league.leagues l WHERE l.id::text = $1`, id)
	if err != nil {
		return League{}, fmt.Errorf("load league: %w", err)
	}
	l, err := pgx.CollectExactlyOneRow(rows, scanLeague)
	if errors.Is(err, pgx.ErrNoRows) {
		return League{}, ErrLeagueNotFound
	}
	return l, err
}

func (LeagueStore) SetRules(ctx context.Context, tx pgx.Tx, id string, divisions Divisions, confirm bool) error {
	_, err := tx.Exec(ctx, `UPDATE league.leagues SET divisions = $2, confirm = $3, updated_at = now() WHERE id::text = $1`, id, divisions, confirm)
	if err != nil {
		return fmt.Errorf("save league rules: %w", err)
	}
	return nil
}

func scanLeague(row pgx.CollectableRow) (League, error) {
	var l League
	err := row.Scan(&l.ID, &l.OriginKind, &l.OriginID, &l.Divisions, &l.Confirm, &l.UpdatedAt)
	return l, err
}
