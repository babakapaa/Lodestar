package auth

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/tenant"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

func NewRoutes(pool *pgxpool.Pool, s Settings) Routes {
	tx := database.NewTenantTx(pool)
	users := identity.NewStore(tx)
	sessions := NewSessionStore(tx, s.Lifetimes)
	login := NewLoginService(LoginDeps{
		Tenants:   tenant.NewStore(pool),
		Users:     users,
		Providers: NewProviderStore(tx, s.Box),
		Attempts:  NewAttemptStore(tx, s.Lifetimes.Login),
		Grants:    NewGrantStore(tx, s.Lifetimes.Grant),
		Sessions:  sessions,
		Discovery: NewDiscovery(s.HTTPClient),
		Redirects: NewRedirects(s.Redirects),
		PublicURL: s.PublicURL,
	})
	return Routes{
		Login:         login,
		Sessions:      sessions,
		Users:         users,
		Authenticator: NewAuthenticator(sessions),
	}
}
