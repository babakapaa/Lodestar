package app

import (
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/core/profile"
	"github.com/babakapaa/lodestar/server/internal/core/schema"
	"github.com/babakapaa/lodestar/server/internal/core/settings"
	"github.com/babakapaa/lodestar/server/internal/core/structure"
	"github.com/babakapaa/lodestar/server/internal/platform/config"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/secretbox"
)

func NewHandler(version string, cfg config.Config, box secretbox.Box, pool *pgxpool.Pool, registry *module.Registry) (http.Handler, error) {
	tx := database.NewTenantTx(pool)
	authRoutes := auth.NewRoutes(pool, auth.Settings{
		Box:        box,
		Lifetimes:  Lifetimes(cfg),
		Redirects:  cfg.AppRedirects,
		PublicURL:  cfg.PublicURL,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	})
	authn := authRoutes.Authenticator

	permissions := append(append(structure.Permissions(), settings.Permissions()...), registry.Permissions()...)
	catalog, err := access.NewCatalog(permissions...)
	if err != nil {
		return nil, err
	}
	preferences, err := schema.NewCatalog(append(profile.CorePreferences(), registry.Preferences()...)...)
	if err != nil {
		return nil, err
	}
	tenantSettings, err := schema.NewCatalog(append(settings.CoreSettings(), registry.Settings()...)...)
	if err != nil {
		return nil, err
	}

	accessService := access.NewService(tx, catalog, structure.Tree{})
	settingsService := settings.NewService(tx, accessService, tenantSettings, registry)
	return NewRouter(version, registry, pool, settings.NewGate(settingsService, authn),
		authRoutes,
		access.Routes{Service: accessService, Authenticator: authn},
		structure.Routes{Service: structure.NewService(tx, accessService), Authenticator: authn},
		profile.Routes{Service: profile.NewService(tx, preferences), Authenticator: authn},
		settings.Routes{Service: settingsService, Authenticator: authn},
	), nil
}
