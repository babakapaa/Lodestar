package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/core/apierror"
	"github.com/babakapaa/lodestar/server/internal/core/health"
	"github.com/babakapaa/lodestar/server/internal/core/meta"
	"github.com/babakapaa/lodestar/server/internal/core/module"
)

func newTestRouter(t *testing.T) http.Handler {
	t.Helper()
	registry, err := module.NewRegistry(Modules()...)
	if err != nil {
		t.Fatalf("NewRegistry: %v", err)
	}
	return NewRouter("test", registry)
}

func serve(t *testing.T, router http.Handler, method, path string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(method, path, nil))
	return rec
}

func TestHealth(t *testing.T) {
	rec := serve(t, newTestRouter(t), http.MethodGet, "/healthz")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var body health.Status
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Status != health.StatusOK {
		t.Fatalf("status field = %q", body.Status)
	}
}

func TestMeta(t *testing.T) {
	rec := serve(t, newTestRouter(t), http.MethodGet, "/v1/meta")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var body meta.Response
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Version != "test" {
		t.Fatalf("version = %q, want test", body.Version)
	}
	if body.Modules == nil {
		t.Fatal("modules must be an empty list, not null")
	}
}

func TestUnknownPathReturnsErrorCode(t *testing.T) {
	rec := serve(t, newTestRouter(t), http.MethodGet, "/v1/missing")

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
	var body apierror.Error
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Code != apierror.CodeNotFound {
		t.Fatalf("code = %q, want %q", body.Code, apierror.CodeNotFound)
	}
}
