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
	return Routes{
		Service: Service{
			tenants:   tenants,
			tx:        tx,
			access:    s.Access,
			links:     s.Links,
			keys:      keys,
			publicURL: s.PublicURL,
		},
		Launch: LaunchService{
			tenants:   tenants,
			tx:        tx,
			users:     identity.NewStore(tx),
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
