package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/babakapaa/lodestar/server/internal/app"
	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/platform/config"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/httpserver"
	"github.com/babakapaa/lodestar/server/internal/platform/secretbox"
)

var version = "dev"

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	if err := run(); err != nil {
		slog.Error("server stopped with error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	box, err := secretbox.New(cfg.SecretKey)
	if err != nil {
		return fmt.Errorf("LODESTAR_SECRET_KEY: %w", err)
	}
	registry, err := module.NewRegistry(app.Modules()...)
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := database.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	migrations := append(app.CoreMigrations(), registry.Migrations()...)
	if err := database.NewMigrator(pool).Apply(ctx, migrations...); err != nil {
		return err
	}

	authRoutes := auth.NewRoutes(pool, auth.Settings{
		Box:        box,
		Lifetimes:  app.Lifetimes(cfg),
		Redirects:  cfg.AppRedirects,
		PublicURL:  cfg.PublicURL,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	})
	catalog, err := access.NewCatalog(registry.Permissions()...)
	if err != nil {
		return err
	}
	accessRoutes := access.Routes{
		Service:       access.NewService(database.NewTenantTx(pool), catalog, access.FlatTree{}),
		Authenticator: authRoutes.Authenticator,
	}
	srv := httpserver.New(cfg.HTTPAddr, app.NewRouter(version, registry, pool, authRoutes, accessRoutes))
	return httpserver.Run(ctx, srv, cfg.ShutdownTimeout)
}
