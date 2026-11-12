package identity

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

const userColumns = `id::text, tenant_id::text, email, display_name, created_at`

type Store struct {
	tenants database.TenantTx
}

func NewStore(tenants database.TenantTx) Store {
	return Store{tenants: tenants}
}

func (s Store) UpsertExternal(ctx context.Context, tenantID string, ext ExternalIdentity) (User, error) {
	var u User
	err := s.tenants.Run(ctx, tenantID, func(tx pgx.Tx) error {
		row := tx.QueryRow(ctx, `
			INSERT INTO core.users (issuer, subject, email, display_name)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (tenant_id, issuer, subject) DO UPDATE
			SET email = EXCLUDED.email, display_name = EXCLUDED.display_name, updated_at = now()
			RETURNING `+userColumns,
			ext.Issuer, ext.Subject, ext.Email, ext.DisplayName)
		var err error
		u, err = scanUser(row)
		return err
	})
	if err != nil {
		return User{}, fmt.Errorf("upsert user: %w", err)
	}
	return u, nil
}

func (s Store) ByID(ctx context.Context, tenantID, userID string) (User, error) {
	var u User
	err := s.tenants.Run(ctx, tenantID, func(tx pgx.Tx) error {
		var err error
		u, err = scanUser(tx.QueryRow(ctx, `SELECT `+userColumns+` FROM core.users WHERE id = $1`, userID))
		return err
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrNotFound
	}
	if err != nil {
		return User{}, fmt.Errorf("find user: %w", err)
	}
	return u, nil
}

func scanUser(row pgx.Row) (User, error) {
	var u User
	err := row.Scan(&u.ID, &u.TenantID, &u.Email, &u.DisplayName, &u.CreatedAt)
	return u, err
}
