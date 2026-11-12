package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/tenant"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

const clientSecretEnv = "LODESTAR_OIDC_CLIENT_SECRET"

func setOIDC(ctx context.Context, args []string) error {
	flags := flag.NewFlagSet("oidc set", flag.ContinueOnError)
	slug := flags.String("tenant", "", "короткое имя вуза")
	key := flags.String("key", "", "короткое имя способа входа, например microsoft")
	name := flags.String("name", "", "название кнопки входа")
	issuer := flags.String("issuer", "", "адрес OIDC провайдера")
	clientID := flags.String("client-id", "", "id приложения у провайдера")
	scopes := flags.String("scopes", "openid,email,profile", "запрашиваемые scopes через запятую")
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

	t, err := tenant.NewStore(e.pool).BySlug(ctx, *slug)
	if err != nil {
		return err
	}
	_, err = auth.NewProviderStore(database.NewTenantTx(e.pool), e.box).Put(ctx, auth.Provider{
		TenantID:     t.ID,
		Key:          *key,
		DisplayName:  *name,
		Issuer:       *issuer,
		ClientID:     *clientID,
		ClientSecret: secret,
		Scopes:       strings.Split(*scopes, ","),
	})
	if err != nil {
		return err
	}
	fmt.Printf("Вход «%s» сохранён для %s.\nУкажите у провайдера адрес возврата:\n%s\n",
		*name, t.Name, auth.CallbackURL(e.cfg.PublicURL, t.Slug, *key))
	return nil
}
