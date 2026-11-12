package structure_test

import (
	"context"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/search"
	"github.com/babakapaa/lodestar/server/internal/core/structure"
	"github.com/babakapaa/lodestar/server/internal/core/tenant"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/database/dbtest"
)

type world struct {
	t         *testing.T
	ctx       context.Context
	tenant    string
	users     identity.Store
	access    access.Service
	structure structure.Service
	admin     auth.Principal
}

func newWorld(t *testing.T) *world {
	t.Helper()
	ctx := context.Background()
	pool := dbtest.New(t)
	sets := []database.Migrations{tenant.Migrations(), search.Migrations(), identity.Migrations(), access.Migrations(), structure.Migrations()}
	if err := database.NewMigrator(pool).Apply(ctx, sets...); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	created, err := tenant.NewStore(pool).Create(ctx, "kaznu", "КазНУ")
	if err != nil {
		t.Fatalf("tenant: %v", err)
	}
	catalog, err := access.NewCatalog(structure.Permissions()...)
	if err != nil {
		t.Fatalf("catalog: %v", err)
	}
	tx := database.NewTenantTx(pool)
	acc := access.NewService(tx, catalog, structure.Tree{})
	w := &world{
		t:         t,
		ctx:       ctx,
		tenant:    created.ID,
		users:     identity.NewStore(tx),
		access:    acc,
		structure: structure.NewService(tx, acc),
	}
	w.admin = w.user("admin", "admin@kaznu.kz")
	if _, err := access.NewBootstrap(tx).GrantTopRole(ctx, w.tenant, w.admin.UserID); err != nil {
		t.Fatalf("grant: %v", err)
	}
	return w
}

func (w *world) user(subject, email string) auth.Principal {
	w.t.Helper()
	u, err := w.users.UpsertExternal(w.ctx, w.tenant, identity.ExternalIdentity{Issuer: "https://idp", Subject: subject, Email: email})
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
