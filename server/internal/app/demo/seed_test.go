package demo_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"slices"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/app"
	"github.com/babakapaa/lodestar/server/internal/app/demo"
	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/structure"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/database/dbtest"
	"github.com/babakapaa/lodestar/server/internal/platform/secretbox"
)

func TestSeedIsRepeatableAndGivesEachPersonTheirRole(t *testing.T) {
	ctx := context.Background()
	pool := dbtest.New(t)
	if err := database.NewMigrator(pool).Apply(ctx, app.CoreMigrations()...); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	box, err := secretbox.New(base64.StdEncoding.EncodeToString(bytes.Repeat([]byte{5}, secretbox.KeySize)))
	if err != nil {
		t.Fatalf("secretbox: %v", err)
	}
	cfg := demo.Config{Issuer: "http://idp.test", ClientID: "lodestar-demo", ClientSecret: "secret", PublicURL: "http://lodestar.test"}

	first, err := demo.Seed(ctx, pool, box, cfg)
	if err != nil {
		t.Fatalf("seed: %v", err)
	}
	second, err := demo.Seed(ctx, pool, box, cfg)
	if err != nil {
		t.Fatalf("seed again: %v", err)
	}
	if first.Tenant.ID != second.Tenant.ID || second.Callback != "http://lodestar.test/v1/auth/demo/demo/callback" {
		t.Fatalf("results = %+v / %+v", first, second)
	}

	tx := database.NewTenantTx(pool)
	catalog, err := access.NewCatalog(app.CorePermissions()...)
	if err != nil {
		t.Fatalf("catalog: %v", err)
	}
	service := access.NewService(tx, catalog, structure.Tree{})
	permissionsOf := func(email string) []access.Permission {
		users, err := identity.NewStore(tx).ByEmail(ctx, first.Tenant.ID, email)
		if err != nil || len(users) != 1 {
			t.Fatalf("user %s: %v %d", email, err, len(users))
		}
		perms, err := service.Permissions(ctx, auth.Principal{TenantID: first.Tenant.ID, UserID: users[0].ID}, "")
		if err != nil {
			t.Fatalf("permissions %s: %v", email, err)
		}
		return perms
	}
	if !slices.Contains(permissionsOf(demo.Admin.Email), access.Administrator) {
		t.Errorf("admin is not an administrator")
	}
	teacher := permissionsOf(demo.Teacher.Email)
	if !slices.Contains(teacher, structure.ManageMembers) || slices.Contains(teacher, access.Administrator) {
		t.Errorf("teacher permissions = %v", teacher)
	}
	if student := permissionsOf(demo.Student.Email); slices.Contains(student, structure.ManageMembers) {
		t.Errorf("student permissions = %v", student)
	}
}
