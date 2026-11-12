package access_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/tenant"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/database/dbtest"
)

const (
	postsWrite access.Permission = "posts.write"
	postsPin   access.Permission = "posts.pin"
)

type nestedTree struct{}

func (nestedTree) Chain(_ context.Context, _ pgx.Tx, tenantID string, s access.Scope) ([]access.Scope, error) {
	chain := []access.Scope{access.TenantScope(tenantID)}
	if s.Kind == "group" {
		chain = append(chain, access.Scope{Kind: "faculty", ID: "f1"})
	}
	if !s.IsTenant() {
		chain = append(chain, s)
	}
	return chain, nil
}

type world struct {
	t       *testing.T
	ctx     context.Context
	service access.Service
	tenant  string
	users   identity.Store
	admin   auth.Principal
}

func newWorld(t *testing.T) *world {
	t.Helper()
	ctx := context.Background()
	pool := dbtest.New(t)
	err := database.NewMigrator(pool).Apply(ctx, tenant.Migrations(), identity.Migrations(), access.Migrations())
	if err != nil {
		t.Fatalf("migrate: %v", err)
	}
	created, err := tenant.NewStore(pool).Create(ctx, "kaznu", "КазНУ")
	if err != nil {
		t.Fatalf("create tenant: %v", err)
	}
	catalog, err := access.NewCatalog(string(postsWrite), string(postsPin))
	if err != nil {
		t.Fatalf("catalog: %v", err)
	}
	tx := database.NewTenantTx(pool)
	w := &world{
		t:       t,
		ctx:     ctx,
		service: access.NewService(tx, catalog, nestedTree{}),
		tenant:  created.ID,
		users:   identity.NewStore(tx),
	}
	w.admin = w.user("admin")
	if _, err := access.NewBootstrap(tx).GrantTopRole(ctx, w.tenant, w.admin.UserID); err != nil {
		t.Fatalf("grant top: %v", err)
	}
	return w
}

func (w *world) user(subject string) auth.Principal {
	w.t.Helper()
	u, err := w.users.UpsertExternal(w.ctx, w.tenant, identity.ExternalIdentity{Issuer: "https://idp", Subject: subject})
	if err != nil {
		w.t.Fatalf("user %s: %v", subject, err)
	}
	return auth.Principal{TenantID: w.tenant, UserID: u.ID}
}

func (w *world) role(name string, perms ...access.Permission) access.Role {
	w.t.Helper()
	r, err := w.service.CreateRole(w.ctx, w.admin, access.RoleInput{Name: name, Permissions: perms})
	if err != nil {
		w.t.Fatalf("create role %s: %v", name, err)
	}
	return r
}

func (w *world) assign(to auth.Principal, r access.Role, scope string) access.Assignment {
	w.t.Helper()
	a, err := w.service.Assign(w.ctx, w.admin, access.AssignInput{RoleID: r.ID, UserID: to.UserID, Scope: scope})
	if err != nil {
		w.t.Fatalf("assign %s: %v", r.Name, err)
	}
	return a
}

func (w *world) can(p auth.Principal, scope string, perm access.Permission) bool {
	w.t.Helper()
	perms, err := w.service.Permissions(w.ctx, p, scope)
	if err != nil {
		w.t.Fatalf("permissions: %v", err)
	}
	return access.NewSet(perms...).Has(perm)
}

func (w *world) defaultRole() access.Role {
	w.t.Helper()
	roles, err := w.service.Roles(w.ctx, w.admin)
	if err != nil {
		w.t.Fatalf("roles: %v", err)
	}
	return roles[len(roles)-1]
}
