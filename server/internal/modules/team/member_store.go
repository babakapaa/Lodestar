package team

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type MemberStore struct{}

func (MemberStore) ByTeams(ctx context.Context, tx pgx.Tx, teamIDs []string) (map[string][]Member, error) {
	out := map[string][]Member{}
	if len(teamIDs) == 0 {
		return out, nil
	}
	rows, err := tx.Query(ctx, `
		SELECT m.team_id::text, u.id::text, u.display_name, u.email, m.label, m.joined_at
		FROM team.members m JOIN core.users u ON u.id = m.user_id
		WHERE m.team_id::text = ANY($1) ORDER BY m.joined_at, u.display_name`, teamIDs)
	if err != nil {
		return nil, fmt.Errorf("team members: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var teamID string
		var m Member
		if err := rows.Scan(&teamID, &m.UserID, &m.DisplayName, &m.Email, &m.Label, &m.JoinedAt); err != nil {
			return nil, fmt.Errorf("team members: %w", err)
		}
		out[teamID] = append(out[teamID], m)
	}
	return out, rows.Err()
}

func (MemberStore) Placed(ctx context.Context, tx pgx.Tx, kind, scopeID string) (map[string]bool, error) {
	rows, err := tx.Query(ctx, `SELECT user_id::text FROM team.members WHERE scope_kind = $1 AND scope_id = $2`, kind, scopeID)
	if err != nil {
		return nil, fmt.Errorf("placed students: %w", err)
	}
	ids, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return nil, fmt.Errorf("placed students: %w", err)
	}
	out := make(map[string]bool, len(ids))
	for _, id := range ids {
		out[id] = true
	}
	return out, nil
}

func (MemberStore) Add(ctx context.Context, tx pgx.Tx, t Team, userID string) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO team.members (team_id, scope_kind, scope_id, user_id) VALUES ($1::uuid, $2, $3, $4::uuid)`,
		t.ID, t.ScopeKind, t.ScopeID, userID)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrInTeam
	}
	if err != nil {
		return fmt.Errorf("add member: %w", err)
	}
	return nil
}

func (MemberStore) Remove(ctx context.Context, tx pgx.Tx, teamID, userID string) error {
	if _, err := tx.Exec(ctx, `DELETE FROM team.members WHERE team_id::text = $1 AND user_id::text = $2`, teamID, userID); err != nil {
		return fmt.Errorf("remove member: %w", err)
	}
	return nil
}

func (MemberStore) SetLabel(ctx context.Context, tx pgx.Tx, teamID, userID, label string) error {
	tag, err := tx.Exec(ctx, `UPDATE team.members SET label = $3 WHERE team_id::text = $1 AND user_id::text = $2`, teamID, userID, label)
	if err != nil {
		return fmt.Errorf("set label: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrTeamNotFound
	}
	return nil
}
