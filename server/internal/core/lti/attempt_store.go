package lti

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type AttemptStore struct {
	ttl time.Duration
}

func (s AttemptStore) Create(ctx context.Context, tx pgx.Tx, a Attempt) error {
	if _, err := tx.Exec(ctx, `DELETE FROM lti.attempts WHERE expires_at < now()`); err != nil {
		return fmt.Errorf("prune lti attempts: %w", err)
	}
	_, err := tx.Exec(ctx, `
		INSERT INTO lti.attempts (state, platform_id, nonce, expires_at)
		VALUES ($1, $2, $3, now() + $4::interval)`, a.State, a.PlatformID, a.Nonce, s.ttl)
	if err != nil {
		return fmt.Errorf("save lti attempt: %w", err)
	}
	return nil
}

func (AttemptStore) Take(ctx context.Context, tx pgx.Tx, state string) (Attempt, error) {
	a := Attempt{State: state}
	err := tx.QueryRow(ctx, `
		DELETE FROM lti.attempts WHERE state = $1 AND expires_at > now()
		RETURNING platform_id::text, nonce`, state).Scan(&a.PlatformID, &a.Nonce)
	if errors.Is(err, pgx.ErrNoRows) {
		return Attempt{}, ErrInvalidState
	}
	if err != nil {
		return Attempt{}, fmt.Errorf("take lti attempt: %w", err)
	}
	return a, nil
}
