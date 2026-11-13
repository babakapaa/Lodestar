package league

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

const grantColumns = `g.id::text, g.award_id::text, a.title, a.gpu_hours, g.user_id::text, g.season_id::text, g.reason, g.granted_at`

type GrantStore struct{}

func (GrantStore) Insert(ctx context.Context, tx pgx.Tx, leagueID, seasonID, granterID string, in GrantInput) (Grant, error) {
	var id string
	err := tx.QueryRow(ctx, `
		INSERT INTO league.grants (award_id, league_id, season_id, user_id, reason, granted_by)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id::text`, in.AwardID, leagueID, seasonID, in.UserID, in.Reason, granterID).Scan(&id)
	if err != nil {
		return Grant{}, fmt.Errorf("grant award: %w", err)
	}
	rows, err := tx.Query(ctx, `SELECT `+grantColumns+` FROM league.grants g JOIN league.awards a ON a.id = g.award_id WHERE g.id::text = $1`, id)
	if err != nil {
		return Grant{}, fmt.Errorf("load grant: %w", err)
	}
	return pgx.CollectExactlyOneRow(rows, scanGrant)
}

func (GrantStore) Delete(ctx context.Context, tx pgx.Tx, id string) error {
	if _, err := tx.Exec(ctx, `DELETE FROM league.grants WHERE id::text = $1`, id); err != nil {
		return fmt.Errorf("revoke award: %w", err)
	}
	return nil
}

func (GrantStore) League(ctx context.Context, tx pgx.Tx, id string) (string, error) {
	var leagueID string
	err := tx.QueryRow(ctx, `SELECT league_id::text FROM league.grants WHERE id::text = $1`, id).Scan(&leagueID)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrGrantNotFound
	}
	if err != nil {
		return "", fmt.Errorf("load grant: %w", err)
	}
	return leagueID, nil
}

func (GrantStore) Of(ctx context.Context, tx pgx.Tx, leagueID string, userID *string) ([]Grant, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+grantColumns+` FROM league.grants g JOIN league.awards a ON a.id = g.award_id
		WHERE g.league_id::text = $1 AND ($2::text IS NULL OR g.user_id::text = $2) ORDER BY g.granted_at, g.id`, leagueID, userID)
	if err != nil {
		return nil, fmt.Errorf("list grants: %w", err)
	}
	return pgx.CollectRows(rows, scanGrant)
}

func scanGrant(row pgx.CollectableRow) (Grant, error) {
	var g Grant
	err := row.Scan(&g.ID, &g.AwardID, &g.Title, &g.GPUHours, &g.UserID, &g.SeasonID, &g.Reason, &g.GrantedAt)
	return g, err
}
