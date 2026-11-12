package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/tenant"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

func grantTopRole(ctx context.Context, args []string) error {
	flags := flag.NewFlagSet("role grant-top", flag.ContinueOnError)
	slug := flags.String("tenant", "", "короткое имя вуза")
	email := flags.String("email", "", "почта человека")
	if err := flags.Parse(args); err != nil {
		return err
	}
	e, err := openEnv(ctx)
	if err != nil {
		return err
	}
	defer e.pool.Close()

	t, err := tenant.NewStore(e.pool).BySlug(ctx, *slug)
	if err != nil {
		return err
	}
	tx := database.NewTenantTx(e.pool)
	users, err := identity.NewStore(tx).ByEmail(ctx, t.ID, *email)
	if err != nil {
		return err
	}
	switch len(users) {
	case 0:
		return fmt.Errorf("%s ещё не входил в %s", *email, t.Name)
	case 1:
	default:
		return fmt.Errorf("с почтой %s несколько учётных записей, уточните в базе", *email)
	}
	role, err := access.NewBootstrap(tx).GrantTopRole(ctx, t.ID, users[0].ID)
	if err != nil {
		return err
	}
	fmt.Printf("%s получил роль «%s» в %s\n", *email, role.Name, t.Name)
	return nil
}
