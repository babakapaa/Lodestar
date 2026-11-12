package app

import (
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/delivery"
	"github.com/babakapaa/lodestar/server/internal/core/files"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/live"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/core/notify"
	"github.com/babakapaa/lodestar/server/internal/core/profile"
	"github.com/babakapaa/lodestar/server/internal/core/schema"
	"github.com/babakapaa/lodestar/server/internal/core/search"
	"github.com/babakapaa/lodestar/server/internal/core/settings"
	"github.com/babakapaa/lodestar/server/internal/core/structure"
	"github.com/babakapaa/lodestar/server/internal/platform/config"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/secretbox"
)

const (
	liveRetention  = 7 * 24 * time.Hour
	pruneEvery     = time.Hour
	uploadDeadline = 24 * time.Hour
	sweepEvery     = time.Hour
	dispatchEvery  = 5 * time.Second
)

type Server struct {
	Handler    http.Handler
	Workers    []Worker
	OnShutdown []func()
}

func NewServer(version string, cfg config.Config, box secretbox.Box, pool *pgxpool.Pool, registry *module.Registry, ext Integrations) (Server, error) {
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
		return Server{}, err
	}
	preferences, err := schema.NewCatalog(append(append(profile.CorePreferences(), delivery.CorePreferences()...), registry.Preferences()...)...)
	if err != nil {
		return Server{}, err
	}
	tenantSettings, err := schema.NewCatalog(append(append(append(settings.CoreSettings(), files.CoreSettings()...), delivery.CoreSettings()...), registry.Settings()...)...)
	if err != nil {
		return Server{}, err
	}
	notifications, err := notify.NewCatalog(registry.NotificationTypes()...)
	if err != nil {
		return Server{}, err
	}
	hub := live.NewHub()

	accessService := access.NewService(tx, catalog, structure.Tree{})
	settingsService := settings.NewService(tx, accessService, tenantSettings, registry)
	profileService := profile.NewService(tx, preferences)
	sources := append([]search.Source{profile.PeopleSource{}, structure.UnitsSource{}}, registry.SearchSources()...)
	searchService, err := search.NewService(tx, settingsService, sources...)
	if err != nil {
		return Server{}, err
	}
	handler := NewRouter(version, registry, pool, settings.NewGate(settingsService, authn),
		authRoutes,
		access.Routes{Service: accessService, Authenticator: authn},
		structure.Routes{Service: structure.NewService(tx, accessService), Authenticator: authn},
		profile.Routes{Service: profileService, Authenticator: authn},
		settings.Routes{Service: settingsService, Authenticator: authn},
		live.Routes{Hub: hub, Tenants: tx, Authenticator: authn, Origins: cfg.WebOrigins},
		notify.Routes{Service: notify.NewService(tx, notifications), Authenticator: authn},
		files.Routes{Service: files.NewService(tx, ext.Storage, settingsService), Authenticator: authn},
		delivery.Routes{Service: delivery.NewService(tx), Authenticator: authn},
		search.Routes{Service: searchService, Authenticator: authn},
	)
	dispatcher := delivery.NewDispatcher(pool, delivery.DispatcherConfig{
		Texts:       delivery.NewTexts(registry.NotificationTypes()...),
		Preferences: profileService,
		Settings:    settingsService,
		Directory:   identity.Store{},
		Pusher:      ext.Pusher,
		Mailer:      ext.Mailer,
		AppURL:      cfg.AppURL,
		Every:       dispatchEvery,
	})
	workers := []Worker{live.NewListener(pool, hub), live.NewPruner(pool, liveRetention, pruneEvery), dispatcher}
	if ext.Storage != nil {
		workers = append(workers, files.NewSweeper(pool, ext.Storage, uploadDeadline, sweepEvery))
	}
	return Server{
		Handler:    handler,
		Workers:    workers,
		OnShutdown: []func(){hub.Close},
	}, nil
}
