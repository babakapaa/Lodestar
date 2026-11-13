package module

import (
	"context"
	"encoding/json"
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

type rosteringModule struct {
	stubModule
	kind string
}

func (m rosteringModule) RosterSources(Kernel) []RosterSource {
	return []RosterSource{rosterKind(m.kind)}
}

type rosterKind string

func (k rosterKind) Kind() string {
	return string(k)
}

func (rosterKind) All(context.Context, pgx.Tx) ([]Roster, error) {
	return nil, nil
}

func (rosterKind) One(context.Context, pgx.Tx, string) (Roster, error) {
	return Roster{}, ErrRosterNotFound
}

func TestRegistryRequiresOwnRosterKinds(t *testing.T) {
	if _, err := NewRegistry(rosteringModule{stubModule{name: "chat"}, "study.course"}); !errors.Is(err, ErrForeignRoster) {
		t.Fatalf("err = %v, want ErrForeignRoster", err)
	}
	registry, err := NewRegistry(rosteringModule{stubModule{name: "study"}, "study.course"})
	if err != nil {
		t.Fatalf("registry: %v", err)
	}
	if _, ok := registry.Rosters(Kernel{}).Find("study.course"); !ok {
		t.Fatal("roster source is not collected")
	}
}

type teamingModule struct {
	stubModule
	kind string
}

func (m teamingModule) TeamScopeSources(Kernel) []TeamScopeSource {
	return []TeamScopeSource{teamScopeKind(m.kind)}
}

type teamScopeKind string

func (k teamScopeKind) Kind() string {
	return string(k)
}

func (teamScopeKind) TeamScope(context.Context, pgx.Tx, string) (TeamScope, error) {
	return TeamScope{}, ErrTeamScopeNotFound
}

func (teamScopeKind) PrepareTeamWork(context.Context, pgx.Tx, string, auth.Principal, json.RawMessage) (json.RawMessage, error) {
	return nil, ErrTeamScopeNotFound
}

func (teamScopeKind) DeliverTeamWork(context.Context, pgx.Tx, TeamDelivery) error {
	return ErrTeamScopeNotFound
}

func TestRegistryRequiresOwnTeamScopeKinds(t *testing.T) {
	if _, err := NewRegistry(teamingModule{stubModule{name: "team"}, "study.assignment"}); !errors.Is(err, ErrForeignTeamScope) {
		t.Fatalf("err = %v, want ErrForeignTeamScope", err)
	}
	registry, err := NewRegistry(teamingModule{stubModule{name: "study"}, "study.assignment"})
	if err != nil {
		t.Fatalf("registry: %v", err)
	}
	scopes := registry.TeamScopes(Kernel{})
	if _, err := scopes.Load(context.Background(), nil, "study.assignment", "x"); !errors.Is(err, ErrTeamScopeNotFound) {
		t.Fatalf("load = %v", err)
	}
	if _, err := scopes.Load(context.Background(), nil, "league.task", "x"); !errors.Is(err, ErrTeamScopeNotFound) {
		t.Fatalf("unknown kind = %v", err)
	}
}

type scoringModule struct {
	stubModule
	kind string
}

func (m scoringModule) ScoreSources(Kernel) []ScoreSource {
	return []ScoreSource{scoreKind(m.kind)}
}

type scoreKind string

func (k scoreKind) Kind() string {
	return string(k)
}

func (scoreKind) ScoredWorks(context.Context, pgx.Tx, string) ([]ScoredWork, error) {
	return nil, nil
}

func TestRegistryRequiresOwnScoreKinds(t *testing.T) {
	if _, err := NewRegistry(scoringModule{stubModule{name: "league"}, "study.course"}); !errors.Is(err, ErrForeignScoreSource) {
		t.Fatalf("err = %v, want ErrForeignScoreSource", err)
	}
	registry, err := NewRegistry(scoringModule{stubModule{name: "study"}, "study.course"})
	if err != nil {
		t.Fatalf("registry: %v", err)
	}
	if _, ok := registry.ScoreSources(Kernel{}).Find("study.course"); !ok {
		t.Fatal("score source is not collected")
	}
}

func TestWorkResultShare(t *testing.T) {
	cases := []struct {
		r    WorkResult
		want float64
	}{
		{WorkResult{Score: 8, Max: 10}, 0.8},
		{WorkResult{Score: 3, Max: 0}, 0},
		{WorkResult{Score: -1, Max: 10}, 0},
	}
	for _, c := range cases {
		if got := c.r.Share(); got != c.want {
			t.Fatalf("share(%v) = %v, want %v", c.r, got, c.want)
		}
	}
}
