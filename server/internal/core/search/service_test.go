package search_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/profile"
	"github.com/babakapaa/lodestar/server/internal/core/search"
	"github.com/babakapaa/lodestar/server/internal/core/structure"
	"github.com/babakapaa/lodestar/server/internal/core/tenant"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/database/dbtest"
)

type modules map[string]bool

func (m modules) ModuleEnabled(_ context.Context, _, name string) (bool, error) {
	return m[name], nil
}

type staticSource struct {
	kind string
	hits []search.Hit
}

func (s staticSource) Kind() string {
	return s.kind
}

func (s staticSource) Search(context.Context, pgx.Tx, auth.Principal, search.Query, int) ([]search.Hit, error) {
	return append([]search.Hit(nil), s.hits...), nil
}

type fixture struct {
	tx database.TenantTx
	me auth.Principal
}

func setup(t *testing.T) fixture {
	t.Helper()
	ctx := context.Background()
	pool := dbtest.New(t)
	sets := []database.Migrations{tenant.Migrations(), search.Migrations(), identity.Migrations(), profile.Migrations(), structure.Migrations()}
	if err := database.NewMigrator(pool).Apply(ctx, sets...); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	created, err := tenant.NewStore(pool).Create(ctx, "kaznu", "КазНУ")
	if err != nil {
		t.Fatalf("tenant: %v", err)
	}
	tx := database.NewTenantTx(pool)
	users := identity.NewStore(tx)
	var me auth.Principal
	for i, name := range []string{"Айгерим Сейтказы", "Бекзат Нурланов", "Семён Ёлкин"} {
		u, err := users.UpsertExternal(ctx, created.ID, identity.ExternalIdentity{Issuer: "https://idp", Subject: string(rune('a' + i)), DisplayName: name})
		if err != nil {
			t.Fatalf("user: %v", err)
		}
		if i == 0 {
			me = auth.Principal{TenantID: created.ID, UserID: u.ID}
		}
	}
	err = tx.Run(ctx, created.ID, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `
			INSERT INTO core.units (kind, name, code) VALUES
			('faculty', 'Факультет физики', 'PHYS'), ('group', 'Физика-21', 'ФЗ-21')`)
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `INSERT INTO core.units (kind, name, archived_at) VALUES ('group', 'Физика-old', now())`)
		return err
	})
	if err != nil {
		t.Fatalf("units: %v", err)
	}
	return fixture{tx: tx, me: me}
}

func (f fixture) service(t *testing.T, extra ...search.Source) search.Service {
	t.Helper()
	sources := append([]search.Source{profile.PeopleSource{}, structure.UnitsSource{}}, extra...)
	svc, err := search.NewService(f.tx, modules{"study": true}, sources...)
	if err != nil {
		t.Fatalf("service: %v", err)
	}
	return svc
}

func TestFindsPeopleByPartOfName(t *testing.T) {
	f := setup(t)
	hits, err := f.service(t).Search(context.Background(), f.me, "бекз", []string{"people"}, 0)
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(hits) != 1 || hits[0].Title != "Бекзат Нурланов" || hits[0].Kind != "people" || hits[0].Link == "" {
		t.Fatalf("hits = %+v", hits)
	}
}

func TestTreatsYoAsYe(t *testing.T) {
	f := setup(t)
	hits, _ := f.service(t).Search(context.Background(), f.me, "Елкин", nil, 0)
	if len(hits) == 0 || hits[0].Title != "Семён Ёлкин" {
		t.Fatalf("hits = %+v", hits)
	}
}

func TestToleratesTypos(t *testing.T) {
	f := setup(t)
	hits, _ := f.service(t).Search(context.Background(), f.me, "Айгирим", []string{"people"}, 0)
	if len(hits) == 0 || hits[0].Title != "Айгерим Сейтказы" {
		t.Fatalf("hits = %+v", hits)
	}
}

func TestFindsUnitsByNameAndCodeButNotArchived(t *testing.T) {
	f := setup(t)
	ctx := context.Background()
	svc := f.service(t)
	hits, _ := svc.Search(ctx, f.me, "физика", []string{"units"}, 0)
	for _, h := range hits {
		if h.Title == "Физика-old" {
			t.Fatalf("archived unit found: %+v", hits)
		}
	}
	if len(hits) == 0 {
		t.Fatal("no units found")
	}
	byCode, _ := svc.Search(ctx, f.me, "фз-21", []string{"units"}, 0)
	if len(byCode) == 0 || byCode[0].Title != "Физика-21" || byCode[0].Subtitle != "group ФЗ-21" {
		t.Fatalf("by code = %+v", byCode)
	}
}

func TestMergesSourcesByScore(t *testing.T) {
	f := setup(t)
	svc := f.service(t,
		staticSource{kind: "study.courses", hits: []search.Hit{{ID: "c1", Title: "Курс", Score: 2}}},
	)
	hits, err := svc.Search(context.Background(), f.me, "физика", nil, 0)
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(hits) < 2 || hits[0].Kind != "study.courses" {
		t.Fatalf("hits = %+v", hits)
	}
}

func TestDisabledModuleIsNotSearched(t *testing.T) {
	f := setup(t)
	svc := f.service(t, staticSource{kind: "notes.pages", hits: []search.Hit{{ID: "n1", Score: 5}}})
	hits, _ := svc.Search(context.Background(), f.me, "физика", nil, 0)
	for _, h := range hits {
		if h.Kind == "notes.pages" {
			t.Fatalf("disabled module answered: %+v", hits)
		}
	}
	kinds, _ := svc.Kinds(context.Background(), f.me)
	if len(kinds) != 2 {
		t.Fatalf("kinds = %v", kinds)
	}
}

func TestRejectsBadQueries(t *testing.T) {
	f := setup(t)
	svc := f.service(t)
	if _, err := svc.Search(context.Background(), f.me, "  ?! ", nil, 0); !errors.Is(err, search.ErrInvalidQuery) {
		t.Fatalf("err = %v", err)
	}
	if _, err := svc.Search(context.Background(), f.me, "физика", []string{"nope"}, 0); !errors.Is(err, search.ErrUnknownKind) {
		t.Fatalf("err = %v", err)
	}
	if _, err := search.NewService(f.tx, modules{}, profile.PeopleSource{}, profile.PeopleSource{}); !errors.Is(err, search.ErrBadSource) {
		t.Fatalf("duplicate source err = %v", err)
	}
}
