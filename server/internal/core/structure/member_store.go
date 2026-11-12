package structure

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type MemberStore struct{}

func (MemberStore) List(ctx context.Context, tx pgx.Tx, unitID string) ([]Member, error) {
	rows, err := tx.Query(ctx, `
		SELECT m.user_id::text, u.display_name, u.email, m.joined_at
		FROM core.unit_members m JOIN core.users u ON u.id = m.user_id
		WHERE m.unit_id::text = $1 ORDER BY u.display_name, u.email`, unitID)
	if err != nil {
		return nil, fmt.Errorf("list members: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (Member, error) {
		var m Member
		err := row.Scan(&m.UserID, &m.DisplayName, &m.Email, &m.JoinedAt)
		return m, err
	})
}

func (MemberStore) Add(ctx context.Context, tx pgx.Tx, unitID, userID string) (bool, error) {
	tag, err := tx.Exec(ctx, `
		INSERT INTO core.unit_members (unit_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, unitID, userID)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && (pgErr.Code == foreignKeyViolation || pgErr.Code == invalidText) {
		return false, ErrUserNotFound
	}
	if err != nil {
		return false, fmt.Errorf("add member: %w", err)
	}
	return tag.RowsAffected() == 1, nil
}

func (MemberStore) Remove(ctx context.Context, tx pgx.Tx, unitID, userID string) (bool, error) {
	tag, err := tx.Exec(ctx, `DELETE FROM core.unit_members WHERE unit_id::text = $1 AND user_id::text = $2`, unitID, userID)
	if err != nil {
		return false, fmt.Errorf("remove member: %w", err)
	}
	return tag.RowsAffected() == 1, nil
}

func (MemberStore) UnitsOf(ctx context.Context, tx pgx.Tx, userID string) ([]Unit, error) {
	rows, err := tx.Query(ctx, `
		SELECT `+prefixed("u")+` FROM core.unit_members m JOIN core.units u ON u.id = m.unit_id
		WHERE m.user_id::text = $1 AND u.archived_at IS NULL ORDER BY u.name`, userID)
	if err != nil {
		return nil, fmt.Errorf("user units: %w", err)
	}
	return pgx.CollectRows(rows, scanUnit)
}

func prefixed(alias string) string {
	return alias + `.id::text, ` + alias + `.parent_id::text, ` + alias + `.kind, ` + alias + `.name, ` +
		alias + `.code, ` + alias + `.archived_at, ` + alias + `.created_at`
}
