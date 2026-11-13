package team

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type InvitationStore struct{}

func (InvitationStore) Put(ctx context.Context, tx pgx.Tx, teamID, userID, invitedBy string) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO team.invitations (team_id, user_id, invited_by) VALUES ($1::uuid, $2::uuid, $3::uuid)
		ON CONFLICT (team_id, user_id) DO NOTHING`, teamID, userID, invitedBy)
	if err != nil {
		return fmt.Errorf("invite: %w", err)
	}
	return nil
}

func (InvitationStore) Delete(ctx context.Context, tx pgx.Tx, teamID, userID string) (bool, error) {
	tag, err := tx.Exec(ctx, `DELETE FROM team.invitations WHERE team_id::text = $1 AND user_id::text = $2`, teamID, userID)
	if err != nil {
		return false, fmt.Errorf("delete invitation: %w", err)
	}
	return tag.RowsAffected() > 0, nil
}

func (InvitationStore) DeleteInScope(ctx context.Context, tx pgx.Tx, kind, scopeID, userID string) error {
	_, err := tx.Exec(ctx, `
		DELETE FROM team.invitations i USING team.teams t
		WHERE t.id = i.team_id AND t.scope_kind = $1 AND t.scope_id = $2 AND i.user_id::text = $3`, kind, scopeID, userID)
	if err != nil {
		return fmt.Errorf("delete invitations: %w", err)
	}
	return nil
}

func (InvitationStore) Invited(ctx context.Context, tx pgx.Tx, teamID string) ([]string, error) {
	rows, err := tx.Query(ctx, `SELECT user_id::text FROM team.invitations WHERE team_id::text = $1 ORDER BY created_at`, teamID)
	if err != nil {
		return nil, fmt.Errorf("invited: %w", err)
	}
	return pgx.CollectRows(rows, pgx.RowTo[string])
}

func (InvitationStore) Mine(ctx context.Context, tx pgx.Tx, userID string) ([]invitationRow, error) {
	rows, err := tx.Query(ctx, `
		SELECT i.team_id::text, coalesce(i.invited_by::text, '') FROM team.invitations i
		JOIN team.teams t ON t.id = i.team_id
		WHERE i.user_id::text = $1 AND t.archived_at IS NULL ORDER BY i.created_at DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("my invitations: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (invitationRow, error) {
		var r invitationRow
		err := row.Scan(&r.teamID, &r.invitedBy)
		return r, err
	})
}
