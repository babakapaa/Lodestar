package lti

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type IdentityStore struct{}

func (IdentityStore) User(ctx context.Context, tx pgx.Tx, platformID, subject string) (string, bool, error) {
	var userID string
	err := tx.QueryRow(ctx, `
		SELECT user_id::text FROM lti.identities WHERE platform_id::text = $1 AND subject = $2`, platformID, subject).Scan(&userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", false, nil
	}
	if err != nil {
		return "", false, fmt.Errorf("find lti identity: %w", err)
	}
	return userID, true, nil
}

func (s IdentityStore) Bind(ctx context.Context, tx pgx.Tx, platformID, subject, userID string) (string, error) {
	_, err := tx.Exec(ctx, `
		INSERT INTO lti.identities (platform_id, subject, user_id) VALUES ($1, $2, $3)
		ON CONFLICT (platform_id, subject) DO NOTHING`, platformID, subject, userID)
	if err != nil {
		return "", fmt.Errorf("bind lti identity: %w", err)
	}
	bound, _, err := s.User(ctx, tx, platformID, subject)
	return bound, err
}
