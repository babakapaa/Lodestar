package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/babakapaa/lodestar/server/internal/app"
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

	handler, err := app.NewHandler(version, cfg, box, pool, registry)
	if err != nil {
		return err
	}
	return httpserver.Run(ctx, httpserver.New(cfg.HTTPAddr, handler), cfg.ShutdownTimeout)
}
