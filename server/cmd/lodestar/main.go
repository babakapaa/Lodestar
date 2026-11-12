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

	ext, err := app.OpenIntegrations(ctx, cfg)
	if err != nil {
		return err
	}
	server, err := app.NewServer(version, cfg, box, pool, registry, ext)
	if err != nil {
		return err
	}
	for _, w := range server.Workers {
		go w.Run(ctx)
	}
	srv := httpserver.New(cfg.HTTPAddr, server.Handler)
	for _, f := range server.OnShutdown {
		srv.RegisterOnShutdown(f)
	}
	return httpserver.Run(ctx, srv, cfg.ShutdownTimeout)
}
