package identity_test

import (
	"context"
	"errors"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/tenant"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/database/dbtest"
)

func TestUsersAreIsolatedByTenant(t *testing.T) {
	ctx := context.Background()
	pool := dbtest.New(t)
	if err := database.NewMigrator(pool).Apply(ctx, tenant.Migrations(), identity.Migrations()); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	tenants := tenant.NewStore(pool)
	a, err := tenants.Create(ctx, "uni-a", "A")
	if err != nil {
		t.Fatalf("create a: %v", err)
	}
	b, err := tenants.Create(ctx, "uni-b", "B")
	if err != nil {
		t.Fatalf("create b: %v", err)
	}
	store := identity.NewStore(database.NewTenantTx(pool))
	ext := identity.ExternalIdentity{Issuer: "https://idp", Subject: "42", Email: "x@a"}

	userA, err := store.UpsertExternal(ctx, a.ID, ext)
	if err != nil {
		t.Fatalf("upsert a: %v", err)
	}
	userB, err := store.UpsertExternal(ctx, b.ID, ext)
	if err != nil {
		t.Fatalf("upsert b: %v", err)
	}
	if userA.ID == userB.ID {
		t.Fatal("two universities share one user record")
	}
	if _, err := store.ByID(ctx, b.ID, userA.ID); !errors.Is(err, identity.ErrNotFound) {
		t.Fatalf("university B sees a user of A: err = %v", err)
	}
}
