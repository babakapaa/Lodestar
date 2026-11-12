package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/babakapaa/lodestar/server/internal/app"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/platform/config"
	"github.com/babakapaa/lodestar/server/internal/platform/httpserver"
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
	registry, err := module.NewRegistry(app.Modules()...)
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv := httpserver.New(cfg.HTTPAddr, app.NewRouter(version, registry))
	return httpserver.Run(ctx, srv, cfg.ShutdownTimeout)
}
