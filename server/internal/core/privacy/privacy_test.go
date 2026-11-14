package privacy_test

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/babakapaa/lodestar/server/internal/core/audit"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/files"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/privacy"
	"github.com/babakapaa/lodestar/server/internal/core/profile"
	"github.com/babakapaa/lodestar/server/internal/core/search"
	"github.com/babakapaa/lodestar/server/internal/core/tenant"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/database/dbtest"
	"github.com/babakapaa/lodestar/server/internal/platform/objectstore"
)

type values map[string]string

func (v values) Value(_ context.Context, _ pgx.Tx, key string) (json.RawMessage, error) {
	return json.RawMessage(v[key]), nil
}

type memoryStorage struct {
	mu      sync.Mutex
	objects map[string]string
}

func (m *memoryStorage) PresignPut(context.Context, string, string, int64, time.Duration) (objectstore.Upload, error) {
	return objectstore.Upload{}, nil
}

func (m *memoryStorage) PresignGet(context.Context, string, string, string, time.Duration) (string, error) {
	return "", nil
}

func (m *memoryStorage) Size(_ context.Context, key string) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	content, ok := m.objects[key]
	if !ok {
		return 0, objectstore.ErrMissing
	}
	return int64(len(content)), nil
}

func (m *memoryStorage) Open(_ context.Context, key string) (io.ReadCloser, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return io.NopCloser(strings.NewReader(m.objects[key])), nil
}

func (m *memoryStorage) Put(_ context.Context, key, _ string, body io.Reader, _ int64) error {
	raw, err := io.ReadAll(body)
	m.mu.Lock()
	defer m.mu.Unlock()
	m.objects[key] = string(raw)
	return err
}

func (m *memoryStorage) Delete(_ context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.objects, key)
	return nil
}

type fixture struct {
	pool     *pgxpool.Pool
	tx       database.TenantTx
	me       auth.Principal
	storage  *memoryStorage
	settings values
	sections privacy.Sections
	file     files.File
}

func setup(t *testing.T) *fixture {
	t.Helper()
	ctx := context.Background()
	pool := dbtest.New(t)
	sets := []database.Migrations{
		tenant.Migrations(), search.Migrations(), identity.Migrations(), audit.Migrations(), privacy.Migrations(),
		profile.Migrations(), files.Migrations(),
	}
	if err := database.NewMigrator(pool).Apply(ctx, sets...); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	created, err := tenant.NewStore(pool).Create(ctx, "kaznu", "КазНУ")
	if err != nil {
		t.Fatalf("tenant: %v", err)
	}
	tx := database.NewTenantTx(pool)
	u, err := identity.NewStore(tx).UpsertExternal(ctx, created.ID, identity.ExternalIdentity{Issuer: "https://idp", Subject: "a", Email: "a@uni.example", DisplayName: "Айгерим"})
	if err != nil {
		t.Fatalf("user: %v", err)
	}
	me := auth.Principal{TenantID: created.ID, UserID: u.ID}
	storage := &memoryStorage{objects: map[string]string{}}
	limits := values{files.MaxFileSetting: `10`}
	fileService := files.NewService(tx, storage, limits)
	upload, err := fileService.Create(ctx, me, files.FileInput{Name: "эссе.txt", Size: 5})
	if err != nil {
		t.Fatalf("create file: %v", err)
	}
	storage.objects[created.ID+"/"+upload.File.ID] = "hello"
	ready, err := fileService.Complete(ctx, me, upload.File.ID)
	if err != nil {
		t.Fatalf("complete file: %v", err)
	}
	sections, err := privacy.NewSections(identity.PrivacySection{}, profile.PrivacySection{}, files.NewPrivacySection(storage), audit.PrivacySection{})
	if err != nil {
		t.Fatalf("sections: %v", err)
	}
	return &fixture{
		pool:     pool,
		tx:       tx,
		me:       me,
		storage:  storage,
		settings: values{privacy.SelfDeletionSetting: `true`, privacy.GraceDaysSetting: `0`},
		sections: sections,
		file:     ready,
	}
}

func TestExportContainsDataAndFiles(t *testing.T) {
	f := setup(t)
	ctx := context.Background()
	archive, err := privacy.NewExporter(f.tx, f.sections).Collect(ctx, f.me)
	if err != nil {
		t.Fatalf("collect: %v", err)
	}
	var buf bytes.Buffer
	if err := archive.Write(ctx, &buf); err != nil {
		t.Fatalf("write: %v", err)
	}
	z, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatalf("zip: %v", err)
	}
	contents := map[string]string{}
	for _, entry := range z.File {
		r, _ := entry.Open()
		b, _ := io.ReadAll(r)
		_ = r.Close()
		contents[entry.Name] = string(b)
	}
	for _, name := range []string{"export.json", "account.json", "profile.json", "files.json", "admin_actions.json"} {
		if _, ok := contents[name]; !ok {
			t.Errorf("%s missing, got %v", name, z.File)
		}
	}
	if !strings.Contains(contents["account.json"], "a@uni.example") {
		t.Errorf("account = %s", contents["account.json"])
	}
	if got := contents["files/"+f.file.ID+"/эссе.txt"]; got != "hello" {
		t.Errorf("file content = %q", got)
	}
}

func TestDeletionCanBeCancelled(t *testing.T) {
	f := setup(t)
	ctx := context.Background()
	f.settings[privacy.GraceDaysSetting] = `14`
	svc := privacy.NewService(f.tx, f.settings)

	d, err := svc.RequestDeletion(ctx, f.me)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	if d.DueAt.Sub(d.RequestedAt) != 14*24*time.Hour {
		t.Fatalf("deletion = %+v", d)
	}
	if err := svc.CancelDeletion(ctx, f.me); err != nil {
		t.Fatalf("cancel: %v", err)
	}
	if _, err := svc.Deletion(ctx, f.me); !errors.Is(err, privacy.ErrNoRequest) {
		t.Fatalf("after cancel err = %v", err)
	}
}

func TestUniversityCanTurnOffSelfDeletion(t *testing.T) {
	f := setup(t)
	f.settings[privacy.SelfDeletionSetting] = `false`
	if _, err := privacy.NewService(f.tx, f.settings).RequestDeletion(context.Background(), f.me); !errors.Is(err, privacy.ErrSelfDeletionOff) {
		t.Fatalf("err = %v", err)
	}
}

func TestDueDeletionErasesEverything(t *testing.T) {
	f := setup(t)
	ctx := context.Background()
	if _, err := privacy.NewService(f.tx, f.settings).RequestDeletion(ctx, f.me); err != nil {
		t.Fatalf("request: %v", err)
	}

	eraser := privacy.NewAccountEraser(f.pool, f.sections, time.Hour).WithClock(func() time.Time { return time.Now().Add(time.Minute) })
	if err := eraser.EraseDue(ctx); err != nil {
		t.Fatalf("erase: %v", err)
	}

	if len(f.storage.objects) != 0 {
		t.Fatalf("objects left: %v", f.storage.objects)
	}
	if _, err := identity.NewStore(f.tx).ByID(ctx, f.me.TenantID, f.me.UserID); !errors.Is(err, identity.ErrNotFound) {
		t.Fatalf("user still exists: %v", err)
	}
	var action string
	err := f.tx.Run(ctx, f.me.TenantID, func(tx pgx.Tx) error {
		return tx.QueryRow(ctx, `SELECT action FROM core.audit_log WHERE target_id = $1`, f.me.UserID).Scan(&action)
	})
	if err != nil || action != privacy.ActionAccountDeleted {
		t.Fatalf("audit = %q, %v", action, err)
	}
}
