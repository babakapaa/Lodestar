package team

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const leaveColumns = `r.id::text, r.team_id::text, t.name, u.id::text, u.display_name, u.email, r.reason, r.status, r.created_at`

type LeaveStore struct{}

func (LeaveStore) Insert(ctx context.Context, tx pgx.Tx, teamID, userID, reason string) (string, error) {
	var id string
	err := tx.QueryRow(ctx, `
		INSERT INTO team.leave_requests (team_id, user_id, reason) VALUES ($1::uuid, $2::uuid, $3) RETURNING id::text`,
		teamID, userID, reason).Scan(&id)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return "", ErrRequestPending
	}
	if err != nil {
		return "", fmt.Errorf("leave request: %w", err)
	}
	return id, nil
}

func (s LeaveStore) ByID(ctx context.Context, tx pgx.Tx, id string) (LeaveRequest, error) {
	list, err := s.many(ctx, tx, `WHERE r.id::text = $1`, id)
	if err != nil {
		return LeaveRequest{}, err
	}
	if len(list) == 0 {
		return LeaveRequest{}, ErrRequestNotFound
	}
	return list[0], nil
}

func (s LeaveStore) Pending(ctx context.Context, tx pgx.Tx, kind, scopeID string) ([]LeaveRequest, error) {
	return s.many(ctx, tx, `WHERE t.scope_kind = $1 AND t.scope_id = $2 AND r.status = 'pending'`, kind, scopeID)
}

func (s LeaveStore) PendingOf(ctx context.Context, tx pgx.Tx, teamID, userID string) (*LeaveRequest, error) {
	list, err := s.many(ctx, tx, `WHERE r.team_id::text = $1 AND r.user_id::text = $2 AND r.status = 'pending'`, teamID, userID)
	if err != nil || len(list) == 0 {
		return nil, err
	}
	return &list[0], nil
}

func (LeaveStore) Decide(ctx context.Context, tx pgx.Tx, id string, status LeaveStatus, decidedBy string) (bool, error) {
	tag, err := tx.Exec(ctx, `
		UPDATE team.leave_requests SET status = $2, decided_by = $3::uuid, decided_at = now()
		WHERE id::text = $1 AND status = 'pending'`, id, status, decidedBy)
	if err != nil {
		return false, fmt.Errorf("decide leave: %w", err)
	}
	return tag.RowsAffected() > 0, nil
}

func (LeaveStore) Close(ctx context.Context, tx pgx.Tx, teamID, userID string) error {
	_, err := tx.Exec(ctx, `
		UPDATE team.leave_requests SET status = 'approved', decided_at = now()
		WHERE team_id::text = $1 AND user_id::text = $2 AND status = 'pending'`, teamID, userID)
	if err != nil {
		return fmt.Errorf("close leave request: %w", err)
	}
	return nil
}

func (LeaveStore) many(ctx context.Context, tx pgx.Tx, where string, args ...any) ([]LeaveRequest, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+leaveColumns+` FROM team.leave_requests r
		JOIN team.teams t ON t.id = r.team_id
		JOIN core.users u ON u.id = r.user_id `+where+` ORDER BY r.created_at`, args...)
	if err != nil {
		return nil, fmt.Errorf("leave requests: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (LeaveRequest, error) {
		var r LeaveRequest
		err := row.Scan(&r.ID, &r.TeamID, &r.TeamName, &r.Student.UserID, &r.Student.DisplayName, &r.Student.Email, &r.Reason, &r.Status, &r.CreatedAt)
		return r, err
	})
}
