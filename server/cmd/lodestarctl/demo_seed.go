package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/babakapaa/lodestar/server/internal/app/demo"
)

func seedDemo(ctx context.Context, args []string) error {
	flags := flag.NewFlagSet("demo seed", flag.ContinueOnError)
	issuer := flags.String("issuer", "http://localhost:8090", "адрес демо-входа lodestar-demo-idp")
	clientID := flags.String("client-id", "lodestar-demo", "id приложения у демо-входа")
	if err := flags.Parse(args); err != nil {
		return err
	}
	secret := os.Getenv(clientSecretEnv)
	if secret == "" {
		return errors.New(clientSecretEnv + " is required")
	}
	e, err := openEnv(ctx)
	if err != nil {
		return err
	}
	defer e.pool.Close()

	result, err := demo.Seed(ctx, e.pool, e.box, demo.Config{
		Issuer:       *issuer,
		ClientID:     *clientID,
		ClientSecret: secret,
		PublicURL:    e.cfg.PublicURL,
	})
	if err != nil {
		return err
	}
	fmt.Printf("Тестовый вуз готов: %s (код %s).\n", result.Tenant.Name, result.Tenant.Slug)
	for _, person := range demo.People() {
		fmt.Printf("  %s, %s: %s\n", person.Name, person.Title, person.Email)
	}
	fmt.Printf("Адрес возврата для DEMO_IDP_REDIRECTS:\n%s\n", result.Callback)
	return nil
}
