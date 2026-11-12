package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/secretbox"
)

const checkViolation = "23514"

type ProviderStore struct {
	tenants database.TenantTx
	box     secretbox.Box
}

func NewProviderStore(tenants database.TenantTx, box secretbox.Box) ProviderStore {
	return ProviderStore{tenants: tenants, box: box}
}

func (s ProviderStore) Put(ctx context.Context, p Provider) (Provider, error) {
	if p.Issuer == "" || p.ClientID == "" || p.ClientSecret == "" {
		return Provider{}, ErrInvalidProvider
	}
	sealed, err := s.box.Seal([]byte(p.ClientSecret))
	if err != nil {
		return Provider{}, err
	}
	if len(p.Scopes) == 0 {
		p.Scopes = []string{"openid", "email", "profile"}
	}
	err = s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		return tx.QueryRow(ctx, `
			INSERT INTO core.identity_providers (key, display_name, issuer, client_id, client_secret, scopes)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (tenant_id, key) DO UPDATE
			SET display_name = EXCLUDED.display_name, issuer = EXCLUDED.issuer,
			    client_id = EXCLUDED.client_id, client_secret = EXCLUDED.client_secret,
			    scopes = EXCLUDED.scopes, updated_at = now()
			RETURNING id::text`,
			p.Key, p.DisplayName, p.Issuer, p.ClientID, sealed, p.Scopes).Scan(&p.ID)
	})
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == checkViolation {
		return Provider{}, ErrInvalidProvider
	}
	if err != nil {
		return Provider{}, fmt.Errorf("save identity provider: %w", err)
	}
	return p, nil
}

func (s ProviderStore) ByKey(ctx context.Context, tenantID, key string) (Provider, error) {
	p := Provider{TenantID: tenantID, Key: key}
	var sealed []byte
	err := s.tenants.Run(ctx, tenantID, func(tx pgx.Tx) error {
		return tx.QueryRow(ctx, `
			SELECT id::text, display_name, issuer, client_id, client_secret, scopes
			FROM core.identity_providers WHERE key = $1`, key).
			Scan(&p.ID, &p.DisplayName, &p.Issuer, &p.ClientID, &sealed, &p.Scopes)
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return Provider{}, ErrProviderNotFound
	}
	if err != nil {
		return Provider{}, fmt.Errorf("find identity provider: %w", err)
	}
	secret, err := s.box.Open(sealed)
	if err != nil {
		return Provider{}, fmt.Errorf("open client secret: %w", err)
	}
	p.ClientSecret = string(secret)
	return p, nil
}

func (s ProviderStore) List(ctx context.Context, tenantID string) ([]ProviderInfo, error) {
	var list []ProviderInfo
	err := s.tenants.Run(ctx, tenantID, func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, `SELECT key, display_name FROM core.identity_providers ORDER BY display_name`)
		if err != nil {
			return err
		}
		list, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (ProviderInfo, error) {
			var info ProviderInfo
			err := row.Scan(&info.Key, &info.DisplayName)
			return info, err
		})
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("list identity providers: %w", err)
	}
	return list, nil
}
