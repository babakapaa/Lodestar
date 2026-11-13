package league_test

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
	"github.com/babakapaa/lodestar/server/internal/modules/league"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/database/dbtest"
)

const course = "course-1"

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

type fakeCourse struct {
	mu     sync.Mutex
	roster module.Roster
	works  map[string]*module.ScoredWork
	order  []string
}

func (*fakeCourse) Kind() string {
	return "study.course"
}

func (f *fakeCourse) All(context.Context, pgx.Tx) ([]module.Roster, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return []module.Roster{f.roster}, nil
}

func (f *fakeCourse) One(_ context.Context, _ pgx.Tx, id string) (module.Roster, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if id != f.roster.ID {
		return module.Roster{}, module.ErrRosterNotFound
	}
	return f.roster, nil
}

func (f *fakeCourse) ScoredWorks(_ context.Context, _ pgx.Tx, id string) ([]module.ScoredWork, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if id != f.roster.ID {
		return nil, nil
	}
	out := make([]module.ScoredWork, 0, len(f.order))
	for _, key := range f.order {
		out = append(out, *f.works[key])
	}
	return out, nil
}

func (f *fakeCourse) work(id, title string, project bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.works[id] = &module.ScoredWork{Kind: "study.assignment", ID: id, Title: title, Project: project}
	f.order = append(f.order, id)
}

func (f *fakeCourse) grade(workID string, p auth.Principal, score, outOf float64) {
	f.mu.Lock()
	defer f.mu.Unlock()
	w := f.works[workID]
	for i, r := range w.Results {
		if r.UserID == p.UserID {
			w.Results[i].Score = score
			return
		}
	}
	w.Results = append(w.Results, module.WorkResult{UserID: p.UserID, Score: score, Max: outOf, At: time.Now()})
}

type allOn struct{}

func (allOn) Value(context.Context, pgx.Tx, string) (json.RawMessage, error) {
	return json.RawMessage(`null`), nil
}

func (allOn) ModuleEnabled(context.Context, string, string) (bool, error) {
	return true, nil
}

type world struct {
	t       *testing.T
	ctx     context.Context
	tenant  string
	tx      database.TenantTx
	users   identity.Store
	leagues league.Service
	worker  league.LeagueWorker
	course  *fakeCourse
	sent    *outbox
	teacher auth.Principal
	aigerim auth.Principal
	erlan   auth.Principal
	dana    auth.Principal
}

func newWorld(t *testing.T) *world {
	t.Helper()
	ctx := context.Background()
	pool := dbtest.New(t)
	sets := []database.Migrations{
		tenant.Migrations(), search.Migrations(), identity.Migrations(), audit.Migrations(), access.Migrations(), structure.Migrations(),
		league.Module{}.Migrations(),
	}
	if err := database.NewMigrator(pool).Apply(ctx, sets...); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	created, err := tenant.NewStore(pool).Create(ctx, "kaznu", "КазНУ")
	if err != nil {
		t.Fatalf("tenant: %v", err)
	}
	tx := database.NewTenantTx(pool)
	fake := &fakeCourse{works: map[string]*module.ScoredWork{}}
	sent := &outbox{}
	k := module.Kernel{
		Pool: pool, Tenants: tx, Notifier: sent, Settings: allOn{},
		Rosters: module.Rosters{fake}, Scores: module.ScoreSources{fake},
	}
	w := &world{
		t: t, ctx: ctx, tenant: created.ID, tx: tx, users: identity.NewStore(tx),
		leagues: league.NewService(k), worker: league.NewLeagueWorker(k), course: fake, sent: sent,
	}
	w.teacher = w.user("teacher", "Преподаватель")
	w.aigerim = w.user("aigerim", "Айгерим")
	w.erlan = w.user("erlan", "Ерлан")
	w.dana = w.user("dana", "Дана")
	fake.roster = module.Roster{
		Kind: "study.course", ID: course, Title: "Основы биологии",
		Leads: []string{w.teacher.UserID}, Members: []string{w.aigerim.UserID, w.erlan.UserID, w.dana.UserID},
	}
	fake.work("project", "Проект", true)
	fake.work("quiz", "Тест", false)
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

func (w *world) view(p auth.Principal) league.LeagueView {
	w.t.Helper()
	v, err := w.leagues.View(w.ctx, p, "study.course", course)
	if err != nil {
		w.t.Fatalf("view: %v", err)
	}
	return v
}

func (w *world) rules(in league.RulesInput) league.LeagueView {
	w.t.Helper()
	v, err := w.leagues.SetRules(w.ctx, w.teacher, w.view(w.teacher).ID, in)
	if err != nil {
		w.t.Fatalf("rules: %v", err)
	}
	return v
}

func (w *world) call(p auth.Principal, method, path, body string) *httptest.ResponseRecorder {
	w.t.Helper()
	mux := http.NewServeMux()
	league.Routes{Leagues: w.leagues}.Mount(mux)
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req = req.WithContext(auth.WithPrincipal(req.Context(), p))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}

func weights(project, quiz int) []league.WeightInput {
	return []league.WeightInput{{Kind: "study.assignment", ID: "project", Weight: project}, {Kind: "study.assignment", ID: "quiz", Weight: quiz}}
}

func three() league.Divisions {
	return league.Divisions{{Name: "Старт", From: 0}, {From: 150}, {Name: "Вершина", From: 300}}
}
