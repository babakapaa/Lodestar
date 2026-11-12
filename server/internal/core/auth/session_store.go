package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

const touchInterval = time.Minute

type SessionStore struct {
	tenants   database.TenantTx
	lifetimes Lifetimes
}

func NewSessionStore(tenants database.TenantTx, lifetimes Lifetimes) SessionStore {
	return SessionStore{tenants: tenants, lifetimes: lifetimes}
}

func (s SessionStore) Start(ctx context.Context, tenantID, userID, deviceName string) (TokenPair, error) {
	access, refresh, err := s.newTokens(tenantID)
	if err != nil {
		return TokenPair{}, err
	}
	err = s.tenants.Run(ctx, tenantID, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `
			INSERT INTO core.sessions
				(user_id, device_name, access_hash, access_expires_at, refresh_hash, refresh_expires_at)
			VALUES ($1, $2, $3, now() + $4::interval, $5, now() + $6::interval)`,
			userID, deviceName, access.Hash, s.lifetimes.Access, refresh.Hash, s.lifetimes.Refresh)
		return err
	})
	if err != nil {
		return TokenPair{}, fmt.Errorf("start session: %w", err)
	}
	return s.pair(access, refresh), nil
}

func (s SessionStore) Refresh(ctx context.Context, value string) (TokenPair, error) {
	old, err := parseToken(refreshPrefix, value)
	if err != nil {
		return TokenPair{}, err
	}
	access, refresh, err := s.newTokens(old.TenantID)
	if err != nil {
		return TokenPair{}, err
	}
	var outcome error
	err = s.tenants.Run(ctx, old.TenantID, func(tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, `
			UPDATE core.sessions
			SET previous_refresh_hash = refresh_hash,
			    refresh_hash = $2, refresh_expires_at = now() + $3::interval,
			    access_hash = $4, access_expires_at = now() + $5::interval,
			    last_used_at = now()
			WHERE refresh_hash = $1 AND revoked_at IS NULL AND refresh_expires_at > now()`,
			old.Hash, refresh.Hash, s.lifetimes.Refresh, access.Hash, s.lifetimes.Access)
		if err != nil || tag.RowsAffected() == 1 {
			return err
		}
		tag, err = tx.Exec(ctx, `
			UPDATE core.sessions SET revoked_at = now()
			WHERE previous_refresh_hash = $1 AND revoked_at IS NULL`, old.Hash)
		if err != nil {
			return err
		}
		outcome = ErrInvalidToken
		if tag.RowsAffected() == 1 {
			outcome = ErrTokenReused
		}
		return nil
	})
	if err != nil {
		return TokenPair{}, fmt.Errorf("refresh session: %w", err)
	}
	if outcome != nil {
		return TokenPair{}, outcome
	}
	return s.pair(access, refresh), nil
}

func (s SessionStore) Authenticate(ctx context.Context, value string) (Principal, error) {
	token, err := parseToken(accessPrefix, value)
	if err != nil {
		return Principal{}, err
	}
	p := Principal{TenantID: token.TenantID}
	err = s.tenants.Run(ctx, token.TenantID, func(tx pgx.Tx) error {
		err := tx.QueryRow(ctx, `
			SELECT id::text, user_id::text FROM core.sessions
			WHERE access_hash = $1 AND revoked_at IS NULL AND access_expires_at > now()`, token.Hash).
			Scan(&p.SessionID, &p.UserID)
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `
			UPDATE core.sessions SET last_used_at = now()
			WHERE id = $1 AND last_used_at < now() - $2::interval`, p.SessionID, touchInterval)
		return err
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return Principal{}, ErrInvalidToken
	}
	if err != nil {
		return Principal{}, fmt.Errorf("authenticate: %w", err)
	}
	return p, nil
}

func (s SessionStore) List(ctx context.Context, p Principal) ([]Session, error) {
	var list []Session
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, `
			SELECT id::text, device_name, created_at, last_used_at FROM core.sessions
			WHERE user_id = $1 AND revoked_at IS NULL AND refresh_expires_at > now()
			ORDER BY last_used_at DESC`, p.UserID)
		if err != nil {
			return err
		}
		list, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (Session, error) {
			var session Session
			err := row.Scan(&session.ID, &session.DeviceName, &session.CreatedAt, &session.LastUsedAt)
			session.Current = session.ID == p.SessionID
			return session, err
		})
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("list sessions: %w", err)
	}
	return list, nil
}

func (s SessionStore) Revoke(ctx context.Context, p Principal, sessionID string) error {
	var revoked int64
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, `
			UPDATE core.sessions SET revoked_at = now()
			WHERE id::text = $1 AND user_id = $2 AND revoked_at IS NULL`, sessionID, p.UserID)
		revoked = tag.RowsAffected()
		return err
	})
	if err != nil {
		return fmt.Errorf("revoke session: %w", err)
	}
	if revoked == 0 {
		return ErrSessionNotFound
	}
	return nil
}

func (s SessionStore) newTokens(tenantID string) (Token, Token, error) {
	access, err := newToken(accessPrefix, tenantID)
	if err != nil {
		return Token{}, Token{}, err
	}
	refresh, err := newToken(refreshPrefix, tenantID)
	if err != nil {
		return Token{}, Token{}, err
	}
	return access, refresh, nil
}

func (s SessionStore) pair(access, refresh Token) TokenPair {
	return TokenPair{
		AccessToken:  access.Value,
		RefreshToken: refresh.Value,
		TokenType:    "Bearer",
		ExpiresIn:    int(s.lifetimes.Access.Seconds()),
	}
}
