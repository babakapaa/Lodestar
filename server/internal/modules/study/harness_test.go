package study_test

import (
	"context"
	"encoding/json"
	"sync"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/audit"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/files"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/core/search"
	"github.com/babakapaa/lodestar/server/internal/core/structure"
	"github.com/babakapaa/lodestar/server/internal/core/tenant"
	"github.com/babakapaa/lodestar/server/internal/modules/study"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/database/dbtest"
)

type world struct {
	t         *testing.T
	ctx       context.Context
	tenant    string
	tx        database.TenantTx
	users     identity.Store
	access    access.Service
	structure structure.Service
	study     study.Service
	work      study.AssignmentService
	files     files.Service
	sent      *outbox
	settings  *fakeSettings
	pool      *pgxpool.Pool
	admin     auth.Principal
}

type fakeSettings struct {
	values map[string]any
}

func (f *fakeSettings) Value(_ context.Context, _ pgx.Tx, key string) (json.RawMessage, error) {
	return json.Marshal(f.values[key])
}

func (f *fakeSettings) ModuleEnabled(context.Context, string, string) (bool, error) {
	return true, nil
}

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

func newWorld(t *testing.T) *world {
	t.Helper()
	ctx := context.Background()
	pool := dbtest.New(t)
	sets := []database.Migrations{
		tenant.Migrations(), search.Migrations(), identity.Migrations(), audit.Migrations(), access.Migrations(),
		structure.Migrations(), files.Migrations(), study.Module{}.Migrations(),
	}
	if err := database.NewMigrator(pool).Apply(ctx, sets...); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	created, err := tenant.NewStore(pool).Create(ctx, "kaznu", "КазНУ")
	if err != nil {
		t.Fatalf("tenant: %v", err)
	}
	catalog, err := access.NewCatalog(append(structure.Permissions(), study.Module{}.Permissions()...)...)
	if err != nil {
		t.Fatalf("catalog: %v", err)
	}
	tx := database.NewTenantTx(pool)
	acc := access.NewService(tx, catalog, structure.Tree{})
	sent := &outbox{}
	w := &world{
		t:         t,
		ctx:       ctx,
		tenant:    created.ID,
		tx:        tx,
		users:     identity.NewStore(tx),
		access:    acc,
		structure: structure.NewService(tx, acc),
		study:     study.NewService(tx, acc),
		work:      study.NewAssignmentService(tx, acc, sent),
		sent:      sent,
		settings:  &fakeSettings{values: map[string]any{study.TimeZoneSetting: "Asia/Almaty", study.TeachersEditSchedule: true}},
		pool:      pool,
	}
	w.files = files.NewService(tx, nil, nil, study.Module{}.FileLinks(module.Kernel{Tenants: tx, Access: acc})...)
	w.admin = w.user("admin", "admin@kaznu.kz")
	if _, err := access.NewBootstrap(tx).GrantTopRole(ctx, w.tenant, w.admin.UserID); err != nil {
		t.Fatalf("grant: %v", err)
	}
	return w
}

func (w *world) user(subject, email string) auth.Principal {
	w.t.Helper()
	u, err := w.users.UpsertExternal(w.ctx, w.tenant, identity.ExternalIdentity{Issuer: "https://idp", Subject: subject, Email: email, DisplayName: subject})
	if err != nil {
		w.t.Fatalf("user: %v", err)
	}
	return auth.Principal{TenantID: w.tenant, UserID: u.ID}
}

func (w *world) unit(parent *structure.Unit, kind, name string) structure.Unit {
	w.t.Helper()
	in := structure.UnitInput{Kind: kind, Name: name}
	if parent != nil {
		in.ParentID = &parent.ID
	}
	u, err := w.structure.Create(w.ctx, w.admin, in)
	if err != nil {
		w.t.Fatalf("create %s: %v", name, err)
	}
	return u
}

func (w *world) join(u structure.Unit, p auth.Principal) {
	w.t.Helper()
	if err := w.structure.AddMember(w.ctx, w.admin, u.ID, p.UserID); err != nil {
		w.t.Fatalf("join: %v", err)
	}
}

func (w *world) course(unit *structure.Unit, title string) study.Course {
	w.t.Helper()
	in := study.CourseInput{Title: title}
	if unit != nil {
		in.UnitID = &unit.ID
	}
	c, err := w.study.Create(w.ctx, w.admin, in)
	if err != nil {
		w.t.Fatalf("course %s: %v", title, err)
	}
	return c
}

func (w *world) add(list study.List, c study.Course, id string) {
	w.t.Helper()
	if err := w.study.Change(w.ctx, w.admin, list, c.ID, id, true); err != nil {
		w.t.Fatalf("add to %s: %v", list, err)
	}
}

func (w *world) grantAt(p auth.Principal, u structure.Unit, perms ...access.Permission) {
	w.t.Helper()
	role, err := w.access.CreateRole(w.ctx, w.admin, access.RoleInput{Name: "r-" + u.Name, Permissions: perms})
	if err != nil {
		w.t.Fatalf("role: %v", err)
	}
	_, err = w.access.Assign(w.ctx, w.admin, access.AssignInput{RoleID: role.ID, UserID: p.UserID, Scope: "unit:" + u.ID})
	if err != nil {
		w.t.Fatalf("assign: %v", err)
	}
}

func (w *world) file(owner auth.Principal, name string) string {
	w.t.Helper()
	var id string
	err := w.tx.Run(w.ctx, w.tenant, func(tx pgx.Tx) error {
		f, err := files.FileStore{}.Create(w.ctx, tx, owner.UserID, files.FileInput{Name: name, ContentType: "text/plain", Size: 1})
		if err != nil {
			return err
		}
		id = f.ID
		_, err = files.FileStore{}.MarkReady(w.ctx, tx, f.ID)
		return err
	})
	if err != nil {
		w.t.Fatalf("file: %v", err)
	}
	return id
}

func (w *world) kernel() module.Kernel {
	return module.Kernel{Pool: w.pool, Tenants: w.tx, Access: w.access, Notifier: w.sent, Settings: w.settings}
}
