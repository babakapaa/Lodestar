package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/setup"
	"github.com/babakapaa/lodestar/server/internal/core/tenant"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

func createTenant(ctx context.Context, args []string) error {
	flags := flag.NewFlagSet("tenant create", flag.ContinueOnError)
	slug := flags.String("slug", "", "короткое имя вуза латиницей, например kaznu")
	name := flags.String("name", "", "полное название вуза")
	ready := flags.Bool("ready", false, "сразу открыть вуз для всех, без мастера установки")
	if err := flags.Parse(args); err != nil {
		return err
	}
	e, err := openEnv(ctx)
	if err != nil {
		return err
	}
	defer e.pool.Close()

	created, err := tenant.NewStore(e.pool).Create(ctx, *slug, *name)
	if err != nil {
		return err
	}
	tx := database.NewTenantTx(e.pool)
	if err := access.NewBootstrap(tx).EnsureDefaults(ctx, created.ID); err != nil {
		return err
	}
	if !*ready {
		err := tx.Run(ctx, created.ID, func(t pgx.Tx) error {
			return setup.Store{}.Begin(ctx, t)
		})
		if err != nil {
			return err
		}
	}
	fmt.Printf("Вуз создан: %s (%s)\n", created.Name, created.Slug)
	return nil
}
