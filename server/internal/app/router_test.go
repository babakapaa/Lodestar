package app

import (
	"context"
	"encoding/json"
	"errors"
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
	return NewRouter("test", registry, stubPinger{}, openGate{})
}

type openGate struct{}

func (openGate) Wrap(_ string, h http.Handler) http.Handler {
	return h
}

type stubPinger struct {
	err error
}

func (p stubPinger) Ping(context.Context) error {
	return p.err
}

func serve(t *testing.T, router http.Handler, path string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
	return rec
}

func TestHealth(t *testing.T) {
	rec := serve(t, newTestRouter(t), "/healthz")

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
	rec := serve(t, newTestRouter(t), "/v1/meta")

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
	rec := serve(t, newTestRouter(t), "/v1/missing")

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

func TestReadyReportsDatabaseDown(t *testing.T) {
	registry, err := module.NewRegistry()
	if err != nil {
		t.Fatalf("NewRegistry: %v", err)
	}
	router := NewRouter("test", registry, stubPinger{err: errors.New("down")}, openGate{})

	rec := serve(t, router, "/readyz")

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want 503", rec.Code)
	}
}

func TestReadyWhenDatabaseUp(t *testing.T) {
	rec := serve(t, newTestRouter(t), "/readyz")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
}
