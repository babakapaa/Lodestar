package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/babakapaa/lodestar/server/internal/core/tenant"
)

func createTenant(ctx context.Context, args []string) error {
	flags := flag.NewFlagSet("tenant create", flag.ContinueOnError)
	slug := flags.String("slug", "", "короткое имя вуза латиницей, например kaznu")
	name := flags.String("name", "", "полное название вуза")
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
	fmt.Printf("Вуз создан: %s (%s)\n", created.Name, created.Slug)
	return nil
}
