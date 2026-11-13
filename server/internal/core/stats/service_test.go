package stats_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/audit"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/core/search"
	"github.com/babakapaa/lodestar/server/internal/core/stats"
	"github.com/babakapaa/lodestar/server/internal/core/structure"
	"github.com/babakapaa/lodestar/server/internal/core/tenant"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/database/dbtest"
)

type headcount struct{}

func (headcount) Stats(_ context.Context, _ pgx.Tx, userIDs []string, from, to time.Time) ([]module.Metric, error) {
	return []module.Metric{{Key: "test.people", Value: float64(len(userIDs))}, {Key: "test.days", Value: to.Sub(from).Hours() / 24}}, nil
}

type world struct {
	t         *testing.T
	ctx       context.Context
	tenant    string
	users     identity.Store
	access    access.Service
	structure structure.Service
	stats     stats.Service
	admin     auth.Principal
	count     int
}

func newWorld(t *testing.T) *world {
	t.Helper()
	ctx := context.Background()
	pool := dbtest.New(t)
	sets := []database.Migrations{tenant.Migrations(), search.Migrations(), identity.Migrations(), audit.Migrations(), access.Migrations(), structure.Migrations()}
	if err := database.NewMigrator(pool).Apply(ctx, sets...); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	created, err := tenant.NewStore(pool).Create(ctx, "kaznu", "КазНУ")
	if err != nil {
		t.Fatalf("tenant: %v", err)
	}
	catalog, err := access.NewCatalog(append(structure.Permissions(), stats.Permissions()...)...)
	if err != nil {
		t.Fatalf("catalog: %v", err)
	}
	tx := database.NewTenantTx(pool)
	acc := access.NewService(tx, catalog, structure.Tree{})
	w := &world{
		t: t, ctx: ctx, tenant: created.ID, users: identity.NewStore(tx), access: acc,
		structure: structure.NewService(tx, acc), stats: stats.NewService(tx, acc, headcount{}),
	}
	w.admin = w.user()
	if _, err := access.NewBootstrap(tx).GrantTopRole(ctx, w.tenant, w.admin.UserID); err != nil {
		t.Fatalf("grant: %v", err)
	}
	return w
}

func (w *world) user() auth.Principal {
	w.t.Helper()
	w.count++
	subject := fmt.Sprintf("u%d", w.count)
	u, err := w.users.UpsertExternal(w.ctx, w.tenant, identity.ExternalIdentity{Issuer: "https://idp", Subject: subject, Email: subject + "@kaznu.kz"})
	if err != nil {
		w.t.Fatalf("user: %v", err)
	}
	return auth.Principal{TenantID: w.tenant, UserID: u.ID}
}

func (w *world) unit(parent *structure.Unit, name string, people int) structure.Unit {
	w.t.Helper()
	in := structure.UnitInput{Kind: "faculty", Name: name}
	if parent != nil {
		in.ParentID = &parent.ID
		in.Kind = "department"
	}
	u, err := w.structure.Create(w.ctx, w.admin, in)
	if err != nil {
		w.t.Fatalf("create %s: %v", name, err)
	}
	for range people {
		if err := w.structure.AddMember(w.ctx, w.admin, u.ID, w.user().UserID); err != nil {
			w.t.Fatalf("member: %v", err)
		}
	}
	return u
}

func (w *world) grant(p auth.Principal, scope string, perms ...access.Permission) {
	w.t.Helper()
	role, err := w.access.CreateRole(w.ctx, w.admin, access.RoleInput{Name: fmt.Sprintf("r%d", w.count), Permissions: perms})
	if err != nil {
		w.t.Fatalf("role: %v", err)
	}
	if _, err := w.access.Assign(w.ctx, w.admin, access.AssignInput{RoleID: role.ID, UserID: p.UserID, Scope: scope}); err != nil {
		w.t.Fatalf("assign: %v", err)
	}
}

func metric(row stats.UnitStats, key string) float64 {
	for _, m := range row.Metrics {
		if m.Key == key {
			return m.Value
		}
	}
	return -1
}

func TestUnitsHideSmallGroups(t *testing.T) {
	w := newWorld(t)
	big := w.unit(nil, "Биофак", 4)
	w.unit(&big, "Кафедра генетики", 2)
	w.unit(nil, "Мехмат", 3)
	rector := w.user()
	w.grant(rector, "tenant", stats.View)

	report, err := w.stats.Units(w.ctx, rector, nil, 0)
	if err != nil {
		t.Fatalf("units: %v", err)
	}
	if report.Days != 30 || report.Threshold != stats.Threshold || report.Full || len(report.Units) != 2 {
		t.Fatalf("report = %+v", report)
	}
	for _, row := range report.Units {
		switch row.Name {
		case "Биофак":
			if row.Hidden || row.People == nil || *row.People != 6 || metric(row, "test.people") != 6 || metric(row, "test.days") != 30 {
				t.Fatalf("big = %+v", row)
			}
		case "Мехмат":
			if !row.Hidden || row.People != nil || len(row.Metrics) != 0 {
				t.Fatalf("small = %+v", row)
			}
		}
	}

	full, err := w.stats.Units(w.ctx, w.admin, nil, 7)
	if err != nil {
		t.Fatalf("admin: %v", err)
	}
	for _, row := range full.Units {
		if row.Hidden || !full.Full {
			t.Fatalf("full = %+v", full)
		}
	}
}

func TestUnitsRespectScope(t *testing.T) {
	w := newWorld(t)
	bio := w.unit(nil, "Биофак", 0)
	w.unit(&bio, "Кафедра генетики", 5)
	dean := w.user()
	w.grant(dean, "unit:"+bio.ID, stats.View)

	if _, err := w.stats.Units(w.ctx, dean, nil, 0); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("tenant err = %v", err)
	}
	report, err := w.stats.Units(w.ctx, dean, &bio.ID, 0)
	if err != nil || len(report.Units) != 1 || report.Units[0].Hidden {
		t.Fatalf("faculty = %+v, %v", report, err)
	}
	if _, err := w.stats.Units(w.ctx, w.user(), &bio.ID, 0); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("stranger err = %v", err)
	}
	if _, err := w.stats.Units(w.ctx, dean, &bio.ID, 401); !errors.Is(err, stats.ErrInvalidPeriod) {
		t.Fatalf("period err = %v", err)
	}
	missing := "00000000-0000-4000-8000-000000000000"
	if _, err := w.stats.Units(w.ctx, dean, &missing, 0); !errors.Is(err, structure.ErrUnitNotFound) {
		t.Fatalf("missing err = %v", err)
	}
}
