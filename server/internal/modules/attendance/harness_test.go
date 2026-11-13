package attendance_test

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
	"github.com/babakapaa/lodestar/server/internal/modules/attendance"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/database/dbtest"
)

type fakeSettings struct {
	mu     sync.Mutex
	values map[string]any
}

func (f *fakeSettings) Value(_ context.Context, _ pgx.Tx, key string) (json.RawMessage, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if v, ok := f.values[key]; ok {
		return json.Marshal(v)
	}
	for _, d := range (attendance.Module{}).Settings() {
		if d.Key == key {
			return json.Marshal(d.Default)
		}
	}
	return json.Marshal("")
}

func (*fakeSettings) ModuleEnabled(context.Context, string, string) (bool, error) {
	return true, nil
}

func (f *fakeSettings) set(key string, value any) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.values[key] = value
}

type fakeLessons struct {
	list []module.LessonRef
}

func (f *fakeLessons) Lesson(_ context.Context, _ pgx.Tx, seriesID, date string) (module.LessonRef, error) {
	for _, l := range f.list {
		if l.SeriesID == seriesID && l.Date == date {
			return l, nil
		}
	}
	return module.LessonRef{}, module.ErrLessonNotFound
}

func (f *fakeLessons) Course(_ context.Context, _ pgx.Tx, courseID, from, to string) (module.LessonCourse, error) {
	out := module.LessonCourse{ID: courseID, Title: "Сети"}
	found := false
	for _, l := range f.list {
		if l.CourseID != courseID {
			continue
		}
		found = true
		out.Leads, out.Members = l.Leads, l.Members
		if l.Date >= from && l.Date <= to {
			out.Lessons = append(out.Lessons, l)
		}
	}
	if !found {
		return out, module.ErrLessonNotFound
	}
	return out, nil
}

type clock struct {
	mu  sync.Mutex
	now time.Time
}

func (c *clock) get() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *clock) move(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = c.now.Add(d)
}

type world struct {
	t        *testing.T
	ctx      context.Context
	tenant   string
	users    identity.Store
	tx       database.TenantTx
	service  attendance.Service
	settings *fakeSettings
	lessons  *fakeLessons
	clock    *clock
	admin    auth.Principal
	teacher  auth.Principal
	aigerim  auth.Principal
	erlan    auth.Principal
	stranger auth.Principal
}

func newWorld(t *testing.T) *world {
	t.Helper()
	ctx := context.Background()
	pool := dbtest.New(t)
	sets := []database.Migrations{
		tenant.Migrations(), search.Migrations(), identity.Migrations(), audit.Migrations(), access.Migrations(), structure.Migrations(),
		attendance.Module{}.Migrations(),
	}
	if err := database.NewMigrator(pool).Apply(ctx, sets...); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	created, err := tenant.NewStore(pool).Create(ctx, "kaznu", "КазНУ")
	if err != nil {
		t.Fatalf("tenant: %v", err)
	}
	tx := database.NewTenantTx(pool)
	catalog, err := access.NewCatalog(append(structure.Permissions(), attendance.Module{}.Permissions()...)...)
	if err != nil {
		t.Fatalf("catalog: %v", err)
	}
	w := &world{
		t: t, ctx: ctx, tenant: created.ID, tx: tx, users: identity.NewStore(tx), settings: &fakeSettings{values: map[string]any{}},
		lessons: &fakeLessons{}, clock: &clock{now: time.Date(2026, 9, 8, 9, 2, 0, 0, time.UTC)},
	}
	w.admin = w.user("admin", "Ректорат")
	if _, err := access.NewBootstrap(tx).GrantTopRole(ctx, w.tenant, w.admin.UserID); err != nil {
		t.Fatalf("grant: %v", err)
	}
	w.teacher = w.user("teacher", "Дана")
	w.aigerim = w.user("aigerim", "Айгерим")
	w.erlan = w.user("erlan", "Ерлан")
	w.stranger = w.user("stranger", "Чужой")
	k := module.Kernel{Pool: pool, Tenants: tx, Settings: w.settings, Access: access.NewService(tx, catalog, structure.Tree{}), Lessons: w.lessons}
	w.service = attendance.NewService(k).WithClock(w.clock.get)
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

func (w *world) lesson(date string, start time.Time) module.LessonRef {
	l := module.LessonRef{
		SeriesID: "s1", CourseID: "course-1", CourseTitle: "Сети", Date: date, StartsAt: start, EndsAt: start.Add(80 * time.Minute),
		Room: "301", Leads: []string{w.teacher.UserID}, Members: []string{w.aigerim.UserID, w.erlan.UserID},
	}
	w.lessons.list = append(w.lessons.list, l)
	return l
}

func (w *world) call(p auth.Principal, method, path, body string) *httptest.ResponseRecorder {
	w.t.Helper()
	mux := http.NewServeMux()
	attendance.Routes{Service: w.service}.Mount(mux)
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req = req.WithContext(auth.WithPrincipal(req.Context(), p))
	req.RemoteAddr = "10.20.30.40:5555"
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}
