package app

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/core/tenant"
	"github.com/babakapaa/lodestar/server/internal/platform/config"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/database/dbtest"
	"github.com/babakapaa/lodestar/server/internal/platform/secretbox"
)

type optionalModule struct{}

func (optionalModule) Manifest() module.Manifest {
	return module.Manifest{Name: "notes", Version: "1", Optional: true}
}

func (optionalModule) Mount(mux *http.ServeMux) {
	mux.HandleFunc("GET /v1/notes/ping", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
}

func TestTurnedOffModuleDisappears(t *testing.T) {
	ctx := context.Background()
	pool := dbtest.New(t)
	if err := database.NewMigrator(pool).Apply(ctx, CoreMigrations()...); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	registry, err := module.NewRegistry(optionalModule{})
	if err != nil {
		t.Fatalf("registry: %v", err)
	}
	box, err := secretbox.New(base64.StdEncoding.EncodeToString(bytes.Repeat([]byte{3}, secretbox.KeySize)))
	if err != nil {
		t.Fatalf("box: %v", err)
	}
	cfg := config.Config{AccessTTL: time.Minute, RefreshTTL: time.Hour, PublicURL: "http://localhost"}
	server, err := NewServer("test", cfg, box, pool, registry)
	if err != nil {
		t.Fatalf("NewServer: %v", err)
	}
	handler := server.Handler

	created, err := tenant.NewStore(pool).Create(ctx, "kaznu", "КазНУ")
	if err != nil {
		t.Fatalf("tenant: %v", err)
	}
	tx := database.NewTenantTx(pool)
	user, err := identity.NewStore(tx).UpsertExternal(ctx, created.ID, identity.ExternalIdentity{Issuer: "https://idp", Subject: "admin"})
	if err != nil {
		t.Fatalf("user: %v", err)
	}
	if _, err := access.NewBootstrap(tx).GrantTopRole(ctx, created.ID, user.ID); err != nil {
		t.Fatalf("grant: %v", err)
	}
	pair, err := auth.NewSessionStore(tx, Lifetimes(cfg)).Start(ctx, created.ID, user.ID, "test")
	if err != nil {
		t.Fatalf("session: %v", err)
	}

	call := func(method, path string, body any) int {
		var data []byte
		if body != nil {
			data, _ = json.Marshal(body)
		}
		req := httptest.NewRequest(method, path, bytes.NewReader(data))
		req.Header.Set("Authorization", "Bearer "+pair.AccessToken)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		return rec.Code
	}

	if status := call(http.MethodGet, "/v1/notes/ping", nil); status != http.StatusNoContent {
		t.Fatalf("enabled module = %d", status)
	}
	if status := call(http.MethodPut, "/v1/modules/notes", map[string]bool{"enabled": false}); status != http.StatusNoContent {
		t.Fatalf("turn off = %d", status)
	}
	if status := call(http.MethodGet, "/v1/notes/ping", nil); status != http.StatusNotFound {
		t.Fatalf("disabled module = %d, want 404", status)
	}
	if status := call(http.MethodPut, "/v1/settings/core.default_theme", map[string]string{"value": `"lab"`}); status != http.StatusBadRequest {
		t.Fatalf("string-wrapped value = %d, want 400", status)
	}
	if status := call(http.MethodPut, "/v1/settings/core.default_theme", map[string]json.RawMessage{"value": json.RawMessage(`"lab"`)}); status != http.StatusNoContent {
		t.Fatalf("set setting = %d", status)
	}
	if status := call(http.MethodGet, "/v1/me/preferences", nil); status != http.StatusOK {
		t.Fatalf("preferences = %d", status)
	}
}
