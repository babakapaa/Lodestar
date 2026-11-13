package team_test

import (
	"context"
	"sync"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/audit"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/core/search"
	"github.com/babakapaa/lodestar/server/internal/core/structure"
	"github.com/babakapaa/lodestar/server/internal/core/tenant"
	"github.com/babakapaa/lodestar/server/internal/modules/team"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/database/dbtest"
)

type outbox struct {
	mu   sync.Mutex
	list []module.Notification
}

func (o *outbox) Notify(_ context.Context, _ pgx.Tx, _ string, n module.Notification) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.list = append(o.list, n)
	return nil
}

func (o *outbox) of(kind string) []module.Notification {
	o.mu.Lock()
	defer o.mu.Unlock()
	var out []module.Notification
	for _, n := range o.list {
		if n.Type == kind {
			out = append(out, n)
		}
	}
	return out
}

type fakeScopes struct {
	mu   sync.Mutex
	list map[string]module.TeamScope
}

func (*fakeScopes) Kind() string {
	return "study.assignment"
}

func (f *fakeScopes) TeamScope(_ context.Context, _ pgx.Tx, id string) (module.TeamScope, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	s, ok := f.list[id]
	if !ok {
		return module.TeamScope{}, module.ErrTeamScopeNotFound
	}
	return s, nil
}

func (f *fakeScopes) set(s module.TeamScope) {
	f.mu.Lock()
	defer f.mu.Unlock()
	s.Kind = "study.assignment"
	f.list[s.ID] = s
}

type world struct {
	t       *testing.T
	ctx     context.Context
	tenant  string
	tx      database.TenantTx
	users   identity.Store
	teams   team.Service
	rosters module.RosterSource
	scopes  *fakeScopes
	sent    *outbox
}

func newWorld(t *testing.T) *world {
	t.Helper()
	ctx := context.Background()
	pool := dbtest.New(t)
	sets := []database.Migrations{
		tenant.Migrations(), search.Migrations(), identity.Migrations(), audit.Migrations(), access.Migrations(), structure.Migrations(),
		team.Module{}.Migrations(),
	}
	if err := database.NewMigrator(pool).Apply(ctx, sets...); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	created, err := tenant.NewStore(pool).Create(ctx, "kaznu", "КазНУ")
	if err != nil {
		t.Fatalf("tenant: %v", err)
	}
	tx := database.NewTenantTx(pool)
	scopes := &fakeScopes{list: map[string]module.TeamScope{}}
	sent := &outbox{}
	k := module.Kernel{Pool: pool, Tenants: tx, Notifier: sent, TeamScopes: module.TeamScopes{scopes}}
	return &world{
		t: t, ctx: ctx, tenant: created.ID, tx: tx, users: identity.NewStore(tx),
		teams: team.NewService(k), rosters: team.Module{}.RosterSources(k)[0], scopes: scopes, sent: sent,
	}
}

func (w *world) user(subject string) auth.Principal {
	w.t.Helper()
	u, err := w.users.UpsertExternal(w.ctx, w.tenant, identity.ExternalIdentity{
		Issuer: "https://idp", Subject: subject, Email: subject + "@kaznu.kz", DisplayName: subject,
	})
	if err != nil {
		w.t.Fatalf("user: %v", err)
	}
	return auth.Principal{TenantID: w.tenant, UserID: u.ID}
}

func (w *world) roster(teamID string) module.Roster {
	w.t.Helper()
	var r module.Roster
	err := w.tx.Run(w.ctx, w.tenant, func(tx pgx.Tx) error {
		var err error
		r, err = w.rosters.One(w.ctx, tx, teamID)
		return err
	})
	if err != nil {
		w.t.Fatalf("roster: %v", err)
	}
	return r
}
