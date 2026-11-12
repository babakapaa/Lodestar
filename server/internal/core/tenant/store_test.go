package tenant_test

import (
	"context"
	"errors"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/core/tenant"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/database/dbtest"
)

func newStore(t *testing.T) tenant.Store {
	t.Helper()
	pool := dbtest.New(t)
	if err := database.NewMigrator(pool).Apply(context.Background(), tenant.Migrations()); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return tenant.NewStore(pool)
}

func TestCreateAndFind(t *testing.T) {
	store := newStore(t)
	ctx := context.Background()

	created, err := store.Create(ctx, "kaznu", "Казахский национальный университет")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	found, err := store.BySlug(ctx, "kaznu")
	if err != nil {
		t.Fatalf("BySlug: %v", err)
	}
	if found.ID != created.ID || found.Name != created.Name {
		t.Fatalf("found %+v, want %+v", found, created)
	}
}

func TestCreateRejectsTakenSlug(t *testing.T) {
	store := newStore(t)
	ctx := context.Background()
	if _, err := store.Create(ctx, "kaznu", "A"); err != nil {
		t.Fatalf("Create: %v", err)
	}

	if _, err := store.Create(ctx, "kaznu", "B"); !errors.Is(err, tenant.ErrSlugTaken) {
		t.Fatalf("err = %v, want ErrSlugTaken", err)
	}
}

func TestCreateRejectsBadSlug(t *testing.T) {
	store := newStore(t)

	if _, err := store.Create(context.Background(), "Bad Slug!", "A"); !errors.Is(err, tenant.ErrInvalid) {
		t.Fatalf("err = %v, want ErrInvalid", err)
	}
}

func TestBySlugNotFound(t *testing.T) {
	store := newStore(t)

	if _, err := store.BySlug(context.Background(), "missing"); !errors.Is(err, tenant.ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}
