package access

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const assignmentColumns = `id::text, role_id::text, user_id::text, scope_kind, scope_id, expires_at, granted_by::text, granted_at`

type AssignmentStore struct{}

func (AssignmentStore) ActiveRoleIDs(ctx context.Context, tx pgx.Tx, userID string, chain []Scope) ([]string, error) {
	kinds, ids := splitChain(chain)
	rows, err := tx.Query(ctx, `
		SELECT DISTINCT role_id::text FROM core.role_assignments
		WHERE user_id::text = $1
		  AND (expires_at IS NULL OR expires_at > now())
		  AND (scope_kind, scope_id) IN (SELECT * FROM unnest($2::text[], $3::text[]))`,
		userID, kinds, ids)
	if err != nil {
		return nil, fmt.Errorf("load role assignments: %w", err)
	}
	return pgx.CollectRows(rows, pgx.RowTo[string])
}

func (AssignmentStore) Upsert(ctx context.Context, tx pgx.Tx, in AssignInput, s Scope, grantedBy *string) (Assignment, error) {
	rows, err := tx.Query(ctx, `
		INSERT INTO core.role_assignments (role_id, user_id, scope_kind, scope_id, expires_at, granted_by)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (role_id, user_id, scope_kind, scope_id) DO UPDATE
		SET expires_at = EXCLUDED.expires_at, granted_by = EXCLUDED.granted_by, granted_at = now()
		RETURNING `+assignmentColumns,
		in.RoleID, in.UserID, s.Kind, s.ID, in.ExpiresAt, grantedBy)
	if err != nil {
		return Assignment{}, assignmentError(err)
	}
	a, err := pgx.CollectExactlyOneRow(rows, scanAssignment)
	return a, assignmentError(err)
}

func (AssignmentStore) ByID(ctx context.Context, tx pgx.Tx, id string) (Assignment, Scope, error) {
	rows, err := tx.Query(ctx, `SELECT `+assignmentColumns+` FROM core.role_assignments WHERE id::text = $1`, id)
	if err != nil {
		return Assignment{}, Scope{}, fmt.Errorf("find role assignment: %w", err)
	}
	var s Scope
	a, err := pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (Assignment, error) {
		var a Assignment
		err := row.Scan(&a.ID, &a.RoleID, &a.UserID, &s.Kind, &s.ID, &a.ExpiresAt, &a.GrantedBy, &a.GrantedAt)
		a.Scope = s.String()
		return a, err
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return Assignment{}, Scope{}, ErrAssignmentNotFound
	}
	return a, s, err
}

func (AssignmentStore) Delete(ctx context.Context, tx pgx.Tx, id string) error {
	_, err := tx.Exec(ctx, `DELETE FROM core.role_assignments WHERE id::text = $1`, id)
	return err
}

func (AssignmentStore) ByUser(ctx context.Context, tx pgx.Tx, userID string) ([]Assignment, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+assignmentColumns+` FROM core.role_assignments
		WHERE user_id::text = $1 AND (expires_at IS NULL OR expires_at > now())
		ORDER BY granted_at`, userID)
	if err != nil {
		return nil, fmt.Errorf("list role assignments: %w", err)
	}
	return pgx.CollectRows(rows, scanAssignment)
}

func scanAssignment(row pgx.CollectableRow) (Assignment, error) {
	var a Assignment
	var s Scope
	err := row.Scan(&a.ID, &a.RoleID, &a.UserID, &s.Kind, &s.ID, &a.ExpiresAt, &a.GrantedBy, &a.GrantedAt)
	a.Scope = s.String()
	return a, err
}

func assignmentError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case foreignKeyViolation, invalidText:
			return ErrUserNotFound
		case checkViolation:
			return ErrInvalidScope
		}
	}
	if err != nil {
		return fmt.Errorf("save role assignment: %w", err)
	}
	return nil
}

func splitChain(chain []Scope) ([]string, []string) {
	kinds := make([]string, len(chain))
	ids := make([]string, len(chain))
	for i, s := range chain {
		kinds[i], ids[i] = s.Kind, s.ID
	}
	return kinds, ids
}
