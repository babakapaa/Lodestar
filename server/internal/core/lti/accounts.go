package lti

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

type Accounts struct {
	tx         database.TenantTx
	users      identity.Store
	identities IdentityStore
}

func (a Accounts) Resolve(ctx context.Context, tenantID string, p Platform, person Person) (string, error) {
	var userID string
	var found bool
	err := a.tx.Run(ctx, tenantID, func(tx pgx.Tx) error {
		var err error
		userID, found, err = a.identities.User(ctx, tx, p.ID, person.Subject)
		return err
	})
	if err != nil || found {
		return userID, err
	}
	if p.MatchEmail && person.Email != "" {
		matches, err := a.users.ByEmail(ctx, tenantID, person.Email)
		if err != nil {
			return "", err
		}
		if len(matches) == 1 {
			userID = matches[0].ID
		}
	}
	if userID == "" {
		u, err := a.users.UpsertExternal(ctx, tenantID, identity.ExternalIdentity{
			Issuer:      p.Issuer,
			Subject:     person.Subject,
			Email:       person.Email,
			DisplayName: person.Name,
		})
		if err != nil {
			return "", err
		}
		userID = u.ID
	}
	err = a.tx.Run(ctx, tenantID, func(tx pgx.Tx) error {
		var err error
		userID, err = a.identities.Bind(ctx, tx, p.ID, person.Subject, userID)
		return err
	})
	return userID, err
}
