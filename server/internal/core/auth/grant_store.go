package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

type GrantStore struct {
	tenants database.TenantTx
	ttl     time.Duration
}

func NewGrantStore(tenants database.TenantTx, ttl time.Duration) GrantStore {
	return GrantStore{tenants: tenants, ttl: ttl}
}

func (s GrantStore) Create(ctx context.Context, g LoginGrant) (string, error) {
	code, err := newToken(grantPrefix, g.TenantID)
	if err != nil {
		return "", err
	}
	err = s.tenants.Run(ctx, g.TenantID, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, `DELETE FROM core.login_grants WHERE expires_at < now()`); err != nil {
			return err
		}
		_, err := tx.Exec(ctx, `
			INSERT INTO core.login_grants (code_hash, user_id, redirect_uri, code_challenge, device_name, expires_at)
			VALUES ($1, $2, $3, $4, $5, now() + $6::interval)`,
			code.Hash, g.UserID, g.RedirectURI, g.CodeChallenge, g.DeviceName, s.ttl)
		return err
	})
	if err != nil {
		return "", fmt.Errorf("save login grant: %w", err)
	}
	return code.Value, nil
}

func (s GrantStore) Take(ctx context.Context, value string) (LoginGrant, error) {
	code, err := parseToken(grantPrefix, value)
	if err != nil {
		return LoginGrant{}, err
	}
	g := LoginGrant{TenantID: code.TenantID}
	err = s.tenants.Run(ctx, code.TenantID, func(tx pgx.Tx) error {
		return tx.QueryRow(ctx, `
			DELETE FROM core.login_grants WHERE code_hash = $1 AND expires_at > now()
			RETURNING user_id::text, redirect_uri, code_challenge, device_name`, code.Hash).
			Scan(&g.UserID, &g.RedirectURI, &g.CodeChallenge, &g.DeviceName)
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return LoginGrant{}, ErrInvalidToken
	}
	if err != nil {
		return LoginGrant{}, fmt.Errorf("take login grant: %w", err)
	}
	return g, nil
}
