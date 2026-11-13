package lti

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type MemberStore struct{}

func (MemberStore) Record(ctx context.Context, tx pgx.Tx, contextID, userID string, m Member) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO lti.members (context_id, user_id, roles, instructor) VALUES ($1, $2, $3, $4)
		ON CONFLICT (context_id, user_id) DO UPDATE
		SET roles = EXCLUDED.roles, instructor = EXCLUDED.instructor, launched_at = now()`,
		contextID, userID, nonNil(m.Roles), m.Instructor)
	if err != nil {
		return fmt.Errorf("record lti member: %w", err)
	}
	return nil
}

func (MemberStore) Find(ctx context.Context, tx pgx.Tx, contextID, userID string) (Member, error) {
	var m Member
	err := tx.QueryRow(ctx, `
		SELECT roles, instructor FROM lti.members WHERE context_id::text = $1 AND user_id::text = $2`,
		contextID, userID).Scan(&m.Roles, &m.Instructor)
	if errors.Is(err, pgx.ErrNoRows) {
		return Member{}, ErrContextNotFound
	}
	if err != nil {
		return Member{}, fmt.Errorf("find lti member: %w", err)
	}
	return m, nil
}
