package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/audit"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/delivery"
	"github.com/babakapaa/lodestar/server/internal/core/files"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/live"
	"github.com/babakapaa/lodestar/server/internal/core/lti"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/core/notify"
	"github.com/babakapaa/lodestar/server/internal/core/privacy"
	"github.com/babakapaa/lodestar/server/internal/core/profile"
	"github.com/babakapaa/lodestar/server/internal/core/schema"
	"github.com/babakapaa/lodestar/server/internal/core/search"
	"github.com/babakapaa/lodestar/server/internal/core/settings"
	"github.com/babakapaa/lodestar/server/internal/core/setup"
	"github.com/babakapaa/lodestar/server/internal/core/stats"
	"github.com/babakapaa/lodestar/server/internal/core/structure"
	"github.com/babakapaa/lodestar/server/internal/platform/clientip"
	"github.com/babakapaa/lodestar/server/internal/platform/config"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/secretbox"
	"github.com/babakapaa/lodestar/server/internal/platform/webapp"
)

const (
	liveRetention    = 7 * 24 * time.Hour
	pruneEvery       = time.Hour
	uploadDeadline   = 24 * time.Hour
	sweepEvery       = time.Hour
	dispatchEvery    = 5 * time.Second
	auditPruneEvery  = 24 * time.Hour
	eraseEvery       = time.Hour
	ltiSyncEvery     = 30 * time.Minute
	setupRemindEvery = time.Hour
)

type Server struct {
	Handler    http.Handler
	Workers    []Worker
	OnShutdown []func()
}

func NewServer(version string, cfg config.Config, box secretbox.Box, pool *pgxpool.Pool, registry *module.Registry, ext Integrations) (Server, error) {
	tx := database.NewTenantTx(pool)
	loginClient := ext.LoginClient
	if loginClient == nil {
		loginClient = &http.Client{Timeout: 10 * time.Second}
	}
	authRoutes := auth.NewRoutes(pool, auth.Settings{
		Box:        box,
		Lifetimes:  Lifetimes(cfg),
		Redirects:  cfg.AppRedirects,
		PublicURL:  cfg.PublicURL,
		HTTPClient: loginClient,
	})
	authn := authRoutes.Authenticator

	core, err := NewCore(pool, registry, box)
	if err != nil {
		return Server{}, err
	}
	preferences, err := schema.NewCatalog(append(append(profile.CorePreferences(), delivery.CorePreferences()...), registry.Preferences()...)...)
	if err != nil {
		return Server{}, err
	}
	hub := live.NewHub()

	accessService := core.Access
	settingsService := core.Settings
	profileService := profile.NewService(tx, preferences)
	kernel := core.Kernel
	kernel.Model = ext.Model
	kernel.Embedder, kernel.Reranker = ext.Embedder, ext.Reranker
	kernel.Calls = ext.Calls
	kernel.Transcriber = ext.Speech
	if ext.Storage != nil {
		kernel.Objects = ext.Storage
	}
	kernel.Preferences = profileService
	kernel.Help = registry.HelpSource(kernel)
	kernel.Briefs = registry.Briefs(kernel)
	kernel.Scores = registry.ScoreSources(kernel)
	kernel.TeamScopes = registry.TeamScopes(kernel)
	kernel.Teams = registry.TeamDirectory(kernel)
	kernel.Rosters = registry.Rosters(kernel)
	kernel.GPU = registry.GPUSources(kernel)
	kernel.Lessons = registry.LessonSource(kernel)
	sources := append([]search.Source{profile.PeopleSource{}, structure.UnitsSource{}}, registry.SearchSources()...)
	searchService, err := search.NewService(tx, settingsService, sources...)
	if err != nil {
		return Server{}, err
	}
	sections, err := privacy.NewSections(append(privacySections(ext), registry.PrivacySections(kernel)...)...)
	if err != nil {
		return Server{}, err
	}
	integrations := lti.Settings{
		Box:           box,
		Access:        accessService,
		Tickets:       authRoutes.Tickets,
		Authenticator: authn,
		Links:         registry.LinkSources(kernel),
		Scores:        kernel.Scores,
		Values:        settingsService,
		HTTPClient:    &http.Client{Timeout: 30 * time.Second},
		AttemptTTL:    Lifetimes(cfg).Login,
		SyncEvery:     ltiSyncEvery,
		PublicURL:     cfg.PublicURL,
		AppURL:        cfg.AppURL,
	}
	setupService := setup.NewService(tx, accessService, settingsService, kernel.Notifier, registry.PilotSources(kernel), registry.ReadinessSources(kernel)...).
		WithHardware(registry.HardwareSources(kernel), probes(cfg, pool, ext)...)
	proxies, err := clientip.Parse(cfg.TrustedProxies)
	if err != nil {
		return Server{}, fmt.Errorf("LODESTAR_TRUSTED_PROXIES: %w", err)
	}
	handler := NewRouter(version, registry, pool, settings.NewGate(settingsService, authn, setupService), kernel,
		authRoutes,
		access.Routes{Service: accessService, Authenticator: authn},
		structure.Routes{Service: structure.NewService(tx, accessService), Authenticator: authn},
		profile.Routes{Service: profileService, Authenticator: authn},
		settings.Routes{Service: settingsService, Authenticator: authn},
		live.Routes{Hub: hub, Tenants: tx, Authenticator: authn, Origins: cfg.WebOrigins},
		notify.Routes{Service: notify.NewService(tx, core.Notifications), Authenticator: authn},
		files.Routes{Service: files.NewService(tx, ext.Storage, settingsService, registry.FileLinks(kernel)...), Authenticator: authn},
		delivery.Routes{Service: delivery.NewService(tx), Authenticator: authn},
		search.Routes{Service: searchService, Authenticator: authn},
		audit.Routes{Service: audit.NewService(tx, accessService), Authenticator: authn},
		setup.Routes{Service: setupService, Authenticator: authn},
		stats.Routes{Service: stats.NewService(tx, accessService, registry.StatsSources(kernel)...), Authenticator: authn},
		lti.NewRoutes(pool, integrations),
		privacy.Routes{
			Service:       privacy.NewService(tx, settingsService),
			Exporter:      privacy.NewExporter(tx, sections),
			Authenticator: authn,
		},
	)
	dispatcher := delivery.NewDispatcher(pool, delivery.DispatcherConfig{
		Texts:       delivery.NewTexts(NotificationTypes(registry)...),
		Preferences: profileService,
		Settings:    settingsService,
		Directory:   identity.Store{},
		Pusher:      ext.Pusher,
		Mailer:      ext.Mailer,
		AppURL:      cfg.AppURL,
		Every:       dispatchEvery,
	})
	workers := []Worker{
		live.NewListener(pool, hub),
		live.NewPruner(pool, liveRetention, pruneEvery),
		dispatcher,
		audit.NewPruner(pool, settingsService, auditPruneEvery),
		privacy.NewAccountEraser(pool, sections, eraseEvery),
		lti.NewSyncWorker(pool, integrations),
		setup.NewReminder(pool, kernel.Notifier, setupRemindEvery),
	}
	for _, w := range registry.Workers(kernel) {
		workers = append(workers, w)
	}
	if ext.Storage != nil {
		workers = append(workers, files.NewSweeper(pool, ext.Storage, uploadDeadline, sweepEvery))
	}
	if cfg.WebDir != "" {
		handler = webapp.New(cfg.WebDir, handler)
	}
	return Server{
		Handler:    proxies.Wrap(handler),
		Workers:    workers,
		OnShutdown: []func(){hub.Close},
	}, nil
}
