package module

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type stubModule struct {
	name string
}

func (s stubModule) Manifest() Manifest {
	return Manifest{Name: s.name, Version: "1"}
}

func (s stubModule) Mount(mux *http.ServeMux) {
	mux.HandleFunc("GET /v1/"+s.name, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
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
	registry.Mount(mux)

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/study", nil))

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", rec.Code)
	}
}
