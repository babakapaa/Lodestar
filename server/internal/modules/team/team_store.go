package team

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

const teamColumns = `t.id::text, t.scope_kind, t.scope_id, t.name, t.created_at, t.archived_at, t.invite_code`

type TeamStore struct {
	members MemberStore
}

func (s TeamStore) Insert(ctx context.Context, tx pgx.Tx, scope Scope, name, createdBy string) (Team, error) {
	rows, err := tx.Query(ctx, `
		INSERT INTO team.teams AS t (scope_kind, scope_id, name, invite_code, created_by)
		VALUES ($1, $2, $3, $4, $5::uuid) RETURNING `+teamColumns,
		scope.Kind, scope.ID, name, rand.Text(), createdBy)
	if err != nil {
		return Team{}, fmt.Errorf("create team: %w", err)
	}
	t, err := pgx.CollectExactlyOneRow(rows, scanTeam)
	if err != nil {
		return Team{}, fmt.Errorf("create team: %w", err)
	}
	t.Members = []Member{}
	return t, nil
}

func (s TeamStore) ByID(ctx context.Context, tx pgx.Tx, id string, lock bool) (Team, error) {
	query := `SELECT ` + teamColumns + ` FROM team.teams t WHERE t.id::text = $1`
	if lock {
		query += ` FOR UPDATE`
	}
	return s.one(ctx, tx, query, id)
}

func (s TeamStore) ByCode(ctx context.Context, tx pgx.Tx, code string) (Team, error) {
	return s.one(ctx, tx, `SELECT `+teamColumns+` FROM team.teams t WHERE t.invite_code = $1 AND t.archived_at IS NULL`, code)
}

func (s TeamStore) Of(ctx context.Context, tx pgx.Tx, kind, scopeID, userID string) (Team, error) {
	return s.one(ctx, tx, `
		SELECT `+teamColumns+` FROM team.teams t
		JOIN team.members m ON m.team_id = t.id
		WHERE m.scope_kind = $1 AND m.scope_id = $2 AND m.user_id::text = $3`, kind, scopeID, userID)
}

func (s TeamStore) ByScope(ctx context.Context, tx pgx.Tx, kind, scopeID string) ([]Team, error) {
	return s.many(ctx, tx, `
		SELECT `+teamColumns+` FROM team.teams t
		WHERE t.scope_kind = $1 AND t.scope_id = $2 AND t.archived_at IS NULL
		ORDER BY t.created_at, t.name`, kind, scopeID)
}

func (s TeamStore) ByUser(ctx context.Context, tx pgx.Tx, userID string) ([]Team, error) {
	return s.many(ctx, tx, `
		SELECT `+teamColumns+` FROM team.teams t
		JOIN team.members m ON m.team_id = t.id
		WHERE m.user_id::text = $1 ORDER BY t.archived_at DESC NULLS FIRST, t.created_at DESC`, userID)
}

func (s TeamStore) All(ctx context.Context, tx pgx.Tx, ids []string) ([]Team, error) {
	return s.many(ctx, tx, `
		SELECT `+teamColumns+` FROM team.teams t
		WHERE $1::text[] IS NULL OR t.id::text = ANY($1) ORDER BY t.id`, ids)
}

func (TeamStore) Rename(ctx context.Context, tx pgx.Tx, id, name string) error {
	if _, err := tx.Exec(ctx, `UPDATE team.teams SET name = $2 WHERE id::text = $1`, id, name); err != nil {
		return fmt.Errorf("rename team: %w", err)
	}
	return nil
}

func (TeamStore) Archive(ctx context.Context, tx pgx.Tx, id string) error {
	if _, err := tx.Exec(ctx, `UPDATE team.teams SET archived_at = coalesce(archived_at, now()) WHERE id::text = $1`, id); err != nil {
		return fmt.Errorf("archive team: %w", err)
	}
	return nil
}

func (s TeamStore) one(ctx context.Context, tx pgx.Tx, query string, args ...any) (Team, error) {
	list, err := s.many(ctx, tx, query, args...)
	if err != nil {
		return Team{}, err
	}
	if len(list) == 0 {
		return Team{}, ErrTeamNotFound
	}
	return list[0], nil
}

func (s TeamStore) many(ctx context.Context, tx pgx.Tx, query string, args ...any) ([]Team, error) {
	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("load teams: %w", err)
	}
	list, err := pgx.CollectRows(rows, scanTeam)
	if err != nil {
		return nil, fmt.Errorf("load teams: %w", err)
	}
	ids := make([]string, 0, len(list))
	for _, t := range list {
		ids = append(ids, t.ID)
	}
	members, err := s.members.ByTeams(ctx, tx, ids)
	if err != nil {
		return nil, err
	}
	for i := range list {
		list[i].Members = members[list[i].ID]
		if list[i].Members == nil {
			list[i].Members = []Member{}
		}
	}
	return list, nil
}

func scanTeam(row pgx.CollectableRow) (Team, error) {
	var t Team
	err := row.Scan(&t.ID, &t.ScopeKind, &t.ScopeID, &t.Name, &t.CreatedAt, &t.ArchivedAt, &t.inviteCode)
	return t, err
}

func notFound(err error) bool {
	return errors.Is(err, ErrTeamNotFound)
}
