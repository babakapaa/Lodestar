package sandbox_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/audit"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/core/search"
	"github.com/babakapaa/lodestar/server/internal/core/structure"
	"github.com/babakapaa/lodestar/server/internal/core/tenant"
	"github.com/babakapaa/lodestar/server/internal/modules/sandbox"
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

type fakeSettings struct {
	mu     sync.Mutex
	values map[string]any
}

func (f *fakeSettings) Value(_ context.Context, _ pgx.Tx, key string) (json.RawMessage, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	v, ok := f.values[key]
	if !ok {
		v = 20
	}
	return json.Marshal(v)
}

func (*fakeSettings) ModuleEnabled(context.Context, string, string) (bool, error) {
	return true, nil
}

func (f *fakeSettings) set(key string, value any) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.values[key] = value
}

type fakeWork struct {
	teacher string
	student string
}

func (fakeWork) Kind() string {
	return "assignment"
}

func (f fakeWork) Brief(_ context.Context, _ pgx.Tx, p auth.Principal, id string) (module.Brief, error) {
	if id != "a1" || (p.UserID != f.teacher && p.UserID != f.student) {
		return module.Brief{}, module.ErrBriefNotFound
	}
	return module.Brief{Kind: "assignment", ID: id, Title: "Нейросеть для МРТ", Teaches: p.UserID == f.teacher, Teachers: []string{f.teacher}}, nil
}

type fakeTeams struct {
	refs []module.TeamRef
}

func (fakeTeams) Placement(context.Context, pgx.Tx, string, string, string) (module.TeamPlacement, error) {
	return module.TeamPlacement{}, nil
}

func (f fakeTeams) Teams(_ context.Context, _ pgx.Tx, userID string) ([]module.TeamRef, error) {
	var out []module.TeamRef
	for _, r := range f.refs {
		for _, m := range r.Members {
			if m == userID {
				out = append(out, r)
			}
		}
	}
	return out, nil
}

type fakeGPU map[string]float64

func (f fakeGPU) GPUGrants(_ context.Context, _ pgx.Tx, userID string, _, _ time.Time) ([]module.GPUGrant, error) {
	if h, ok := f[userID]; ok {
		return []module.GPUGrant{{Source: "league.division", Title: "Серебро", Hours: h}}, nil
	}
	return nil, nil
}

type world struct {
	t        *testing.T
	ctx      context.Context
	tenant   string
	users    identity.Store
	quotas   sandbox.Service
	sent     *outbox
	settings *fakeSettings
	admin    auth.Principal
	teacher  auth.Principal
	aigerim  auth.Principal
	erlan    auth.Principal
}

func newWorld(t *testing.T) *world {
	t.Helper()
	ctx := context.Background()
	pool := dbtest.New(t)
	sets := []database.Migrations{
		tenant.Migrations(), search.Migrations(), identity.Migrations(), audit.Migrations(), access.Migrations(), structure.Migrations(),
		sandbox.Module{}.Migrations(),
	}
	if err := database.NewMigrator(pool).Apply(ctx, sets...); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	created, err := tenant.NewStore(pool).Create(ctx, "kaznu", "КазНУ")
	if err != nil {
		t.Fatalf("tenant: %v", err)
	}
	tx := database.NewTenantTx(pool)
	catalog, err := access.NewCatalog(append(structure.Permissions(), sandbox.Module{}.Permissions()...)...)
	if err != nil {
		t.Fatalf("catalog: %v", err)
	}
	w := &world{t: t, ctx: ctx, tenant: created.ID, users: identity.NewStore(tx), sent: &outbox{}, settings: &fakeSettings{values: map[string]any{}}}
	w.admin = w.user("admin", "Ректорат")
	if _, err := access.NewBootstrap(tx).GrantTopRole(ctx, w.tenant, w.admin.UserID); err != nil {
		t.Fatalf("grant: %v", err)
	}
	w.teacher = w.user("teacher", "Преподаватель")
	w.aigerim = w.user("aigerim", "Айгерим")
	w.erlan = w.user("erlan", "Ерлан")
	k := module.Kernel{
		Pool: pool, Tenants: tx, Notifier: w.sent, Settings: w.settings, Access: access.NewService(tx, catalog, structure.Tree{}),
		Briefs: module.Briefs{fakeWork{teacher: w.teacher.UserID, student: w.aigerim.UserID}},
		Teams:  fakeTeams{refs: []module.TeamRef{{ID: "team-1", Name: "Нейроны", Members: []string{w.aigerim.UserID, w.erlan.UserID}}}},
		GPU:    module.GPUSources{fakeGPU{w.aigerim.UserID: 10}},
	}
	w.quotas = sandbox.NewService(k)
	return w
}

func (w *world) user(subject, name string) auth.Principal {
	w.t.Helper()
	u, err := w.users.UpsertExternal(w.ctx, w.tenant, identity.ExternalIdentity{
		Issuer: "https://idp", Subject: subject, Email: subject + "@kaznu.kz", DisplayName: name,
	})
	if err != nil {
		w.t.Fatalf("user: %v", err)
	}
	return auth.Principal{TenantID: w.tenant, UserID: u.ID}
}

func (w *world) term() sandbox.Term {
	w.t.Helper()
	today := time.Now().UTC()
	t, err := w.quotas.CreateTerm(w.ctx, w.admin, sandbox.TermInput{
		Name: "Осень", StartsOn: today.AddDate(0, -1, 0).Format(time.DateOnly), EndsOn: today.AddDate(0, 2, 0).Format(time.DateOnly),
	})
	if err != nil {
		w.t.Fatalf("term: %v", err)
	}
	return t
}

func (w *world) quota(p auth.Principal) sandbox.Quota {
	w.t.Helper()
	q, err := w.quotas.Quota(w.ctx, p)
	if err != nil {
		w.t.Fatalf("quota: %v", err)
	}
	return q
}

func (w *world) call(p auth.Principal, method, path, body string) *httptest.ResponseRecorder {
	w.t.Helper()
	mux := http.NewServeMux()
	sandbox.Routes{Quotas: w.quotas}.Mount(mux)
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req = req.WithContext(auth.WithPrincipal(req.Context(), p))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}
