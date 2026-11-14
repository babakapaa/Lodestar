package tester

import (
	"context"
	"crypto/rand"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/babakapaa/lodestar/server/internal/app"
	"github.com/babakapaa/lodestar/server/internal/app/demo"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/platform/config"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/demoidp"
	"github.com/babakapaa/lodestar/server/internal/platform/secretbox"
)

const (
	IdPPath      = "/tester/idp"
	ClientID     = "lodestar-tester"
	loginTimeout = 10 * time.Second
)

type Mode struct {
	tenantID string
	people   map[string]string
	sessions auth.SessionStore
	idp      http.Handler
	issuer   string
}

func Start(ctx context.Context, cfg config.Config, pool *pgxpool.Pool, box secretbox.Box, registry *module.Registry) (Mode, error) {
	issuer := cfg.PublicURL + IdPPath
	secret := rand.Text()
	provider, err := demoidp.New(demoidp.Config{
		Issuer:       issuer,
		ClientID:     ClientID,
		ClientSecret: secret,
		Redirects:    []string{auth.CallbackURL(cfg.PublicURL, demo.TenantSlug, demo.ProviderKey)},
		Users:        demo.People(),
	})
	if err != nil {
		return Mode{}, err
	}
	result, err := demo.Seed(ctx, pool, box, registry, demo.Config{
		Issuer:       issuer,
		ClientID:     ClientID,
		ClientSecret: secret,
		PublicURL:    cfg.PublicURL,
	})
	if err != nil {
		return Mode{}, err
	}
	return Mode{
		tenantID: result.Tenant.ID,
		people:   result.People,
		sessions: auth.NewSessionStore(database.NewTenantTx(pool), app.Lifetimes(cfg)),
		idp:      http.StripPrefix(IdPPath, provider.Handler()),
		issuer:   issuer,
	}, nil
}

func (m Mode) LoginClient() *http.Client {
	return &http.Client{
		Timeout:   loginTimeout,
		Transport: localTransport{prefix: m.issuer + "/", handler: m.idp, next: http.DefaultTransport},
	}
}

func (m Mode) Wrap(next http.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /v1/tester", InfoHandler{})
	mux.Handle("POST /v1/tester/sign-in", SignInHandler{tenantID: m.tenantID, people: m.people, sessions: m.sessions})
	mux.Handle(IdPPath+"/", m.idp)
	mux.Handle("/", next)
	return mux
}
