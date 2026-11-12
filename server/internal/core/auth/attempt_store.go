package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

type AttemptStore struct {
	tenants database.TenantTx
	ttl     time.Duration
}

func NewAttemptStore(tenants database.TenantTx, ttl time.Duration) AttemptStore {
	return AttemptStore{tenants: tenants, ttl: ttl}
}

func (s AttemptStore) Create(ctx context.Context, tenantID string, a LoginAttempt) error {
	err := s.tenants.Run(ctx, tenantID, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, `DELETE FROM core.login_attempts WHERE expires_at < now()`); err != nil {
			return err
		}
		_, err := tx.Exec(ctx, `
			INSERT INTO core.login_attempts
				(state, provider_id, nonce, verifier, redirect_uri, code_challenge, device_name, expires_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, now() + $8::interval)`,
			a.State, a.ProviderID, a.Nonce, a.Verifier, a.RedirectURI, a.CodeChallenge, a.DeviceName, s.ttl)
		return err
	})
	if err != nil {
		return fmt.Errorf("save login attempt: %w", err)
	}
	return nil
}

func (s AttemptStore) Take(ctx context.Context, tenantID, state string) (LoginAttempt, error) {
	a := LoginAttempt{State: state}
	err := s.tenants.Run(ctx, tenantID, func(tx pgx.Tx) error {
		return tx.QueryRow(ctx, `
			DELETE FROM core.login_attempts WHERE state = $1 AND expires_at > now()
			RETURNING provider_id::text, nonce, verifier, redirect_uri, code_challenge, device_name`, state).
			Scan(&a.ProviderID, &a.Nonce, &a.Verifier, &a.RedirectURI, &a.CodeChallenge, &a.DeviceName)
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return LoginAttempt{}, ErrInvalidState
	}
	if err != nil {
		return LoginAttempt{}, fmt.Errorf("take login attempt: %w", err)
	}
	return a, nil
}
