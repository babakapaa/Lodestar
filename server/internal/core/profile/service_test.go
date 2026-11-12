package profile_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/profile"
	"github.com/babakapaa/lodestar/server/internal/core/schema"
	"github.com/babakapaa/lodestar/server/internal/core/search"
	"github.com/babakapaa/lodestar/server/internal/core/tenant"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/database/dbtest"
)

func setup(t *testing.T) (profile.Service, auth.Principal, auth.Principal) {
	t.Helper()
	ctx := context.Background()
	pool := dbtest.New(t)
	if err := database.NewMigrator(pool).Apply(ctx, tenant.Migrations(), search.Migrations(), identity.Migrations(), profile.Migrations()); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	created, err := tenant.NewStore(pool).Create(ctx, "kaznu", "КазНУ")
	if err != nil {
		t.Fatalf("tenant: %v", err)
	}
	tx := database.NewTenantTx(pool)
	users := identity.NewStore(tx)
	principal := func(subject, name string) auth.Principal {
		u, err := users.UpsertExternal(ctx, created.ID, identity.ExternalIdentity{Issuer: "https://idp", Subject: subject, DisplayName: name})
		if err != nil {
			t.Fatalf("user: %v", err)
		}
		return auth.Principal{TenantID: created.ID, UserID: u.ID}
	}
	catalog, err := schema.NewCatalog(profile.CorePreferences()...)
	if err != nil {
		t.Fatalf("catalog: %v", err)
	}
	return profile.NewService(tx, catalog), principal("a", "Айгерим"), principal("b", "Бекзат")
}

func TestProfileDefaultsToAccountName(t *testing.T) {
	svc, me, _ := setup(t)

	p, err := svc.Me(context.Background(), me)
	if err != nil {
		t.Fatalf("me: %v", err)
	}
	if p.DisplayName != "Айгерим" || p.Contacts == nil {
		t.Fatalf("profile = %+v", p)
	}
}

func TestOthersSeeOnlyPublicContacts(t *testing.T) {
	svc, me, other := setup(t)
	ctx := context.Background()
	_, err := svc.Update(ctx, me, profile.ProfileInput{
		DisplayName: "Айка",
		About:       "Люблю физику",
		Contacts: []profile.Contact{
			{Kind: "telegram", Value: "@aika", Public: true},
			{Kind: "phone", Value: "+77000000000"},
		},
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}

	mine, _ := svc.Me(ctx, me)
	seen, err := svc.Of(ctx, other, me.UserID)
	if err != nil {
		t.Fatalf("of: %v", err)
	}
	if mine.DisplayName != "Айка" || len(mine.Contacts) != 2 {
		t.Fatalf("mine = %+v", mine)
	}
	if len(seen.Contacts) != 1 || seen.Contacts[0].Kind != "telegram" {
		t.Fatalf("seen = %+v", seen.Contacts)
	}
}

func TestProfileRejectsBadContact(t *testing.T) {
	svc, me, _ := setup(t)

	_, err := svc.Update(context.Background(), me, profile.ProfileInput{Contacts: []profile.Contact{{Kind: "Bad Kind", Value: "x"}}})
	if !errors.Is(err, profile.ErrInvalidProfile) {
		t.Fatalf("err = %v, want ErrInvalidProfile", err)
	}
}

func TestPreferences(t *testing.T) {
	svc, me, _ := setup(t)
	ctx := context.Background()

	if err := svc.SetPreference(ctx, me, "core.theme", json.RawMessage(`"lab"`)); err != nil {
		t.Fatalf("set: %v", err)
	}
	if err := svc.SetPreference(ctx, me, "core.theme", json.RawMessage(`"neon"`)); !errors.Is(err, schema.ErrInvalidValue) {
		t.Fatalf("bad value: err = %v", err)
	}
	if err := svc.SetPreference(ctx, me, "core.unknown", json.RawMessage(`1`)); !errors.Is(err, schema.ErrUnknownKey) {
		t.Fatalf("unknown key: err = %v", err)
	}
	values := prefs(t, svc, me)
	if string(values["core.theme"].Value) != `"lab"` || !values["core.theme"].Set || values["core.appearance"].Set {
		t.Fatalf("values = %+v", values)
	}
	if err := svc.ResetPreference(ctx, me, "core.theme"); err != nil {
		t.Fatalf("reset: %v", err)
	}
	if values := prefs(t, svc, me); values["core.theme"].Set {
		t.Fatal("reset did not clear the preference")
	}
}

func prefs(t *testing.T, svc profile.Service, p auth.Principal) map[string]schema.Value {
	t.Helper()
	list, err := svc.Preferences(context.Background(), p)
	if err != nil {
		t.Fatalf("preferences: %v", err)
	}
	m := map[string]schema.Value{}
	for _, v := range list {
		m[v.Key] = v
	}
	return m
}
