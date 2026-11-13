package lti

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/tenant"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

func NewRoutes(pool *pgxpool.Pool, s Settings) Routes {
	tx := database.NewTenantTx(pool)
	tenants := tenant.NewStore(pool)
	keys := NewKeyStore(s.Box)
	accounts := Accounts{tx: tx, users: identity.NewStore(tx)}
	return Routes{
		Service: Service{
			tenants:   tenants,
			tx:        tx,
			access:    s.Access,
			links:     s.Links,
			keys:      keys,
			values:    s.Values,
			sync:      newSyncService(tx, s, keys, accounts),
			publicURL: s.PublicURL,
		},
		Launch: LaunchService{
			tenants:   tenants,
			tx:        tx,
			accounts:  accounts,
			tickets:   s.Tickets,
			keySets:   NewKeySets(s.HTTPClient),
			links:     s.Links,
			attempts:  AttemptStore{ttl: s.AttemptTTL},
			publicURL: s.PublicURL,
			appURL:    s.AppURL,
		},
		JWKS:          JWKSHandler{tenants: tenants, tx: tx, keys: keys},
		Authenticator: s.Authenticator,
	}
}

func NewSyncWorker(pool *pgxpool.Pool, s Settings) SyncWorker {
	tx := database.NewTenantTx(pool)
	accounts := Accounts{tx: tx, users: identity.NewStore(tx)}
	return SyncWorker{pool: pool, sync: newSyncService(tx, s, NewKeyStore(s.Box), accounts), every: s.SyncEvery}
}

func newSyncService(tx database.TenantTx, s Settings, keys KeyStore, accounts Accounts) SyncService {
	return SyncService{
		tx:       tx,
		values:   s.Values,
		links:    s.Links,
		scores:   s.Scores,
		client:   NewPlatformClient(s.HTTPClient),
		accounts: accounts,
		keys:     keys,
	}
}
