package module

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/search"
)

type stubModule struct {
	name string
}

func (s stubModule) Manifest() Manifest {
	return Manifest{Name: s.name, Version: "1"}
}

func (s stubModule) Mount(mux *http.ServeMux, _ Kernel) {
	mux.HandleFunc("GET /v1/"+s.name+"/ping", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("GET /v1/"+s.name, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	})
}

func TestRegistryRejectsDuplicates(t *testing.T) {
	_, err := NewRegistry(stubModule{name: "study"}, stubModule{name: "study"})
	if !errors.Is(err, ErrDuplicateModule) {
		t.Fatalf("err = %v, want ErrDuplicateModule", err)
	}
}

func TestRegistryRejectsUnnamed(t *testing.T) {
	_, err := NewRegistry(stubModule{})
	if !errors.Is(err, ErrUnnamedModule) {
		t.Fatalf("err = %v, want ErrUnnamedModule", err)
	}
}

func TestRegistryManifestsKeepOrder(t *testing.T) {
	registry, err := NewRegistry(stubModule{name: "study"}, stubModule{name: "notes"})
	if err != nil {
		t.Fatalf("NewRegistry: %v", err)
	}
	got := registry.Manifests()
	if len(got) != 2 || got[0].Name != "study" || got[1].Name != "notes" {
		t.Fatalf("Manifests = %+v", got)
	}
}

func TestRegistryMountsRoutes(t *testing.T) {
	registry, err := NewRegistry(stubModule{name: "study"})
	if err != nil {
		t.Fatalf("NewRegistry: %v", err)
	}
	mux := http.NewServeMux()
	registry.Mount(mux, closedGate{closed: "notes"}, Kernel{})

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/study/ping", nil))

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", rec.Code)
	}
}

type closedGate struct {
	closed string
}

func (g closedGate) Wrap(name string, h http.Handler) http.Handler {
	if name == g.closed {
		return http.NotFoundHandler()
	}
	return h
}

func TestRegistryMountsTheModuleRoot(t *testing.T) {
	registry, err := NewRegistry(stubModule{name: "notes"})
	if err != nil {
		t.Fatalf("NewRegistry: %v", err)
	}
	mux := http.NewServeMux()
	registry.Mount(mux, closedGate{}, Kernel{})

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/notes", nil))

	if rec.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want 202", rec.Code)
	}
}

func TestRegistryMountGoesThroughGate(t *testing.T) {
	registry, err := NewRegistry(stubModule{name: "notes"})
	if err != nil {
		t.Fatalf("NewRegistry: %v", err)
	}
	mux := http.NewServeMux()
	registry.Mount(mux, closedGate{closed: "notes"}, Kernel{})

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/notes/ping", nil))

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

type permissionedModule struct {
	stubModule
}

func (permissionedModule) Permissions() []string {
	return []string{"study.grades.view"}
}

func TestRegistryCollectsPermissions(t *testing.T) {
	registry, err := NewRegistry(stubModule{name: "notes"}, permissionedModule{stubModule{name: "study"}})
	if err != nil {
		t.Fatalf("NewRegistry: %v", err)
	}
	if got := registry.Permissions(); len(got) != 1 || got[0] != "study.grades.view" {
		t.Fatalf("Permissions = %v", got)
	}
}

func TestRegistryRejectsReservedNames(t *testing.T) {
	for _, name := range []string{"roles", "Study", "x"} {
		if _, err := NewRegistry(stubModule{name: name}); !errors.Is(err, ErrReservedName) {
			t.Errorf("NewRegistry(%q) err = %v", name, err)
		}
	}
}

type searchableModule struct {
	stubModule
	kind string
}

func (m searchableModule) SearchSources() []search.Source {
	return []search.Source{kindOnly(m.kind)}
}

type kindOnly string

func (k kindOnly) Kind() string {
	return string(k)
}

func (kindOnly) Search(context.Context, pgx.Tx, auth.Principal, search.Query, int) ([]search.Hit, error) {
	return nil, nil
}

func TestRegistryRequiresOwnSearchKinds(t *testing.T) {
	if _, err := NewRegistry(searchableModule{stubModule{name: "study"}, "notes.pages"}); !errors.Is(err, ErrForeignSearchKind) {
		t.Fatalf("err = %v, want ErrForeignSearchKind", err)
	}
	registry, err := NewRegistry(searchableModule{stubModule{name: "study"}, "study.courses"})
	if err != nil || len(registry.SearchSources()) != 1 {
		t.Fatalf("registry = %v, %v", registry, err)
	}
}
