package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/babakapaa/lodestar/server/internal/app"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/platform/config"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/secretbox"
)

type env struct {
	cfg      config.Config
	pool     *pgxpool.Pool
	box      secretbox.Box
	registry *module.Registry
}

func openEnv(ctx context.Context) (env, error) {
	cfg, err := config.Load()
	if err != nil {
		return env{}, err
	}
	box, err := secretbox.New(cfg.SecretKey)
	if err != nil {
		return env{}, fmt.Errorf("LODESTAR_SECRET_KEY: %w", err)
	}
	pool, err := database.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		return env{}, err
	}
	registry, err := module.NewRegistry(app.Modules()...)
	if err != nil {
		pool.Close()
		return env{}, err
	}
	if err := database.NewMigrator(pool).Apply(ctx, append(app.CoreMigrations(), registry.Migrations()...)...); err != nil {
		pool.Close()
		return env{}, err
	}
	return env{cfg: cfg, pool: pool, box: box, registry: registry}, nil
}
