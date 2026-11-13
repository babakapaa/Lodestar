package team_test

import (
	"context"
	"encoding/json"
	"net/http"
	"slices"
	"strings"
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
	"github.com/babakapaa/lodestar/server/internal/platform/llm"
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
	mu        sync.Mutex
	list      map[string]module.TeamScope
	delivered []module.TeamDelivery
}

func (f *fakeScopes) PrepareTeamWork(_ context.Context, _ pgx.Tx, _ string, _ auth.Principal, work json.RawMessage) (json.RawMessage, error) {
	if strings.Contains(string(work), "incomplete") {
		return nil, module.TeamWorkProblem{Status: http.StatusUnprocessableEntity, Code: "search_map_incomplete"}
	}
	return work, nil
}

func (f *fakeScopes) DeliverTeamWork(_ context.Context, _ pgx.Tx, d module.TeamDelivery) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.delivered = append(f.delivered, d)
	return nil
}

func (f *fakeScopes) deliveries() []module.TeamDelivery {
	f.mu.Lock()
	defer f.mu.Unlock()
	return slices.Clone(f.delivered)
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

type scriptedModel struct {
	mu      sync.Mutex
	reply   string
	request llm.Request
}

func (m *scriptedModel) Complete(_ context.Context, req llm.Request) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.request = req
	return m.reply, nil
}

func (m *scriptedModel) Stream(ctx context.Context, req llm.Request, delta func(string) error) (string, error) {
	full, err := m.Complete(ctx, req)
	if err != nil {
		return "", err
	}
	return full, delta(full)
}

func (m *scriptedModel) asked() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.request.Messages[len(m.request.Messages)-1].Content
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
	model   *scriptedModel
	kernel  module.Kernel
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
	model := &scriptedModel{}
	k := module.Kernel{Pool: pool, Tenants: tx, Notifier: sent, Settings: allOn{}, TeamScopes: module.TeamScopes{scopes}, Model: model}
	return &world{
		t: t, ctx: ctx, tenant: created.ID, tx: tx, users: identity.NewStore(tx),
		teams: team.NewService(k), rosters: team.Module{}.RosterSources(k)[0], scopes: scopes, sent: sent,
		model: model, kernel: k,
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

type allOn struct{}

func (allOn) Value(context.Context, pgx.Tx, string) (json.RawMessage, error) {
	return json.RawMessage(`null`), nil
}

func (allOn) ModuleEnabled(context.Context, string, string) (bool, error) {
	return true, nil
}

func (w *world) placement(scopeID, userID string) (module.TeamPlacement, error) {
	var out module.TeamPlacement
	err := w.tx.Run(w.ctx, w.tenant, func(tx pgx.Tx) error {
		var err error
		out, err = team.Module{}.TeamDirectory(w.kernel).Placement(w.ctx, tx, "study.assignment", scopeID, userID)
		return err
	})
	return out, err
}
