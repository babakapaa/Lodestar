package files_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/files"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/tenant"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/database/dbtest"
)

type fixture struct {
	pool    *pgxpool.Pool
	storage *memoryStorage
	service files.Service
	me      auth.Principal
	other   auth.Principal
}

func setup(t *testing.T) fixture {
	t.Helper()
	ctx := context.Background()
	pool := dbtest.New(t)
	if err := database.NewMigrator(pool).Apply(ctx, tenant.Migrations(), identity.Migrations(), files.Migrations()); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	created, err := tenant.NewStore(pool).Create(ctx, "kaznu", "КазНУ")
	if err != nil {
		t.Fatalf("tenant: %v", err)
	}
	tx := database.NewTenantTx(pool)
	users := identity.NewStore(tx)
	principal := func(subject string) auth.Principal {
		u, err := users.UpsertExternal(ctx, created.ID, identity.ExternalIdentity{Issuer: "https://idp", Subject: subject})
		if err != nil {
			t.Fatalf("user: %v", err)
		}
		return auth.Principal{TenantID: created.ID, UserID: u.ID}
	}
	storage := newMemoryStorage()
	return fixture{
		pool:    pool,
		storage: storage,
		service: files.NewService(tx, storage, fixedLimit{mb: 1}),
		me:      principal("a"),
		other:   principal("b"),
	}
}

func (f fixture) upload(t *testing.T, size int64) files.File {
	t.Helper()
	ctx := context.Background()
	created, err := f.service.Create(ctx, f.me, files.FileInput{Name: "Курсовая.pdf", ContentType: "application/pdf", Size: size})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	f.storage.put(f.me.TenantID+"/"+created.File.ID, size)
	ready, err := f.service.Complete(ctx, f.me, created.File.ID)
	if err != nil {
		t.Fatalf("complete: %v", err)
	}
	return ready
}

func TestUploadBecomesReadyAfterContentArrives(t *testing.T) {
	f := setup(t)
	ctx := context.Background()
	created, err := f.service.Create(ctx, f.me, files.FileInput{Name: "notes.txt", Size: 10})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.File.Status != files.StatusPending || created.File.ContentType != "application/octet-stream" || created.Upload.Method != "PUT" {
		t.Fatalf("created = %+v", created)
	}
	if _, err := f.service.Complete(ctx, f.me, created.File.ID); !errors.Is(err, files.ErrNotUploaded) {
		t.Fatalf("complete before upload err = %v", err)
	}

	f.storage.put(f.me.TenantID+"/"+created.File.ID, 10)
	ready, err := f.service.Complete(ctx, f.me, created.File.ID)
	if err != nil || ready.Status != files.StatusReady || ready.UploadedAt == nil {
		t.Fatalf("ready = %+v, %v", ready, err)
	}
	list, _ := f.service.List(ctx, f.me, nil)
	if len(list) != 1 {
		t.Fatalf("list = %+v", list)
	}
}

func TestSizeMismatchIsRejected(t *testing.T) {
	f := setup(t)
	ctx := context.Background()
	created, _ := f.service.Create(ctx, f.me, files.FileInput{Name: "a.bin", Size: 10})
	f.storage.put(f.me.TenantID+"/"+created.File.ID, 11)

	if _, err := f.service.Complete(ctx, f.me, created.File.ID); !errors.Is(err, files.ErrSizeMismatch) {
		t.Fatalf("err = %v, want ErrSizeMismatch", err)
	}
}

func TestUniversityLimitApplies(t *testing.T) {
	f := setup(t)
	_, err := f.service.Create(context.Background(), f.me, files.FileInput{Name: "big.mp4", Size: 2 << 20})
	if !errors.Is(err, files.ErrTooLarge) {
		t.Fatalf("err = %v, want ErrTooLarge", err)
	}
}

func TestInvalidInputIsRejected(t *testing.T) {
	f := setup(t)
	for _, in := range []files.FileInput{
		{Name: "", Size: 1},
		{Name: "../etc/passwd", Size: 1},
		{Name: "a.txt", Size: 0},
		{Name: "a.txt", ContentType: "not a type", Size: 1},
	} {
		if _, err := f.service.Create(context.Background(), f.me, in); !errors.Is(err, files.ErrInvalidFile) {
			t.Errorf("Create(%+v) err = %v", in, err)
		}
	}
}

func TestOthersCannotSeeOrDeleteMyFile(t *testing.T) {
	f := setup(t)
	ctx := context.Background()
	mine := f.upload(t, 5)

	if _, err := f.service.Get(ctx, f.other, mine.ID); !errors.Is(err, files.ErrNotFound) {
		t.Fatalf("get err = %v", err)
	}
	if _, err := f.service.Download(ctx, f.other, mine.ID); !errors.Is(err, files.ErrNotFound) {
		t.Fatalf("download err = %v", err)
	}
	if err := f.service.Delete(ctx, f.other, mine.ID); !errors.Is(err, files.ErrNotFound) {
		t.Fatalf("delete err = %v", err)
	}
	if _, err := f.service.Get(ctx, f.me, "not-a-uuid"); !errors.Is(err, files.ErrNotFound) {
		t.Fatalf("malformed id err = %v", err)
	}
}

func TestDeleteRemovesContent(t *testing.T) {
	f := setup(t)
	ctx := context.Background()
	mine := f.upload(t, 5)

	d, err := f.service.Download(ctx, f.me, mine.ID)
	if err != nil || d.URL == "" {
		t.Fatalf("download = %+v, %v", d, err)
	}
	if err := f.service.Delete(ctx, f.me, mine.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if f.storage.has(f.me.TenantID + "/" + mine.ID) {
		t.Fatal("object is still stored")
	}
	if _, err := f.service.Get(ctx, f.me, mine.ID); !errors.Is(err, files.ErrNotFound) {
		t.Fatalf("get after delete err = %v", err)
	}
}

func TestSweeperDropsAbandonedUploads(t *testing.T) {
	f := setup(t)
	ctx := context.Background()
	pending, _ := f.service.Create(ctx, f.me, files.FileInput{Name: "half.zip", Size: 3})
	f.storage.put(f.me.TenantID+"/"+pending.File.ID, 1)
	ready := f.upload(t, 4)

	if err := files.NewSweeper(f.pool, f.storage, 0, time.Hour).SweepAll(ctx); err != nil {
		t.Fatalf("sweep: %v", err)
	}
	if _, err := f.service.Get(ctx, f.me, pending.File.ID); !errors.Is(err, files.ErrNotFound) {
		t.Fatalf("pending file survived: %v", err)
	}
	if f.storage.has(f.me.TenantID + "/" + pending.File.ID) {
		t.Fatal("pending object survived")
	}
	if _, err := f.service.Get(ctx, f.me, ready.ID); err != nil {
		t.Fatalf("ready file was swept: %v", err)
	}
}

func TestWithoutStorageFilesAreUnavailable(t *testing.T) {
	f := setup(t)
	svc := files.NewService(database.NewTenantTx(f.pool), nil, fixedLimit{mb: 1})
	if _, err := svc.Create(context.Background(), f.me, files.FileInput{Name: "a", Size: 1}); !errors.Is(err, files.ErrUnavailable) {
		t.Fatalf("err = %v, want ErrUnavailable", err)
	}
}
