package access_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/tenant"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/database/dbtest"
)

func TestRoutesEnforcePermissions(t *testing.T) {
	ctx := context.Background()
	pool := dbtest.New(t)
	err := database.NewMigrator(pool).Apply(ctx, tenant.Migrations(), identity.Migrations(), auth.Migrations(), access.Migrations())
	if err != nil {
		t.Fatalf("migrate: %v", err)
	}
	created, err := tenant.NewStore(pool).Create(ctx, "kaznu", "КазНУ")
	if err != nil {
		t.Fatalf("tenant: %v", err)
	}
	tx := database.NewTenantTx(pool)
	users := identity.NewStore(tx)
	sessions := auth.NewSessionStore(tx, auth.Lifetimes{Access: time.Minute, Refresh: time.Hour})
	token := func(subject string) (string, string) {
		u, err := users.UpsertExternal(ctx, created.ID, identity.ExternalIdentity{Issuer: "https://idp", Subject: subject})
		if err != nil {
			t.Fatalf("user: %v", err)
		}
		pair, err := sessions.Start(ctx, created.ID, u.ID, "test")
		if err != nil {
			t.Fatalf("session: %v", err)
		}
		return u.ID, pair.AccessToken
	}
	adminID, adminToken := token("admin")
	_, studentToken := token("student")
	if _, err := access.NewBootstrap(tx).GrantTopRole(ctx, created.ID, adminID); err != nil {
		t.Fatalf("grant: %v", err)
	}
	catalog, err := access.NewCatalog()
	if err != nil {
		t.Fatalf("catalog: %v", err)
	}
	mux := http.NewServeMux()
	access.Routes{
		Service:       access.NewService(tx, catalog, access.FlatTree{}),
		Authenticator: auth.NewAuthenticator(sessions),
	}.Mount(mux)
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	call := func(method, path, bearer string, body any) int {
		var data []byte
		if body != nil {
			data, _ = json.Marshal(body)
		}
		req, err := http.NewRequest(method, server.URL+path, bytes.NewReader(data))
		if err != nil {
			t.Fatalf("request: %v", err)
		}
		req.Header.Set("Authorization", "Bearer "+bearer)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("%s %s: %v", method, path, err)
		}
		defer func() { _ = resp.Body.Close() }()
		return resp.StatusCode
	}

	if status := call(http.MethodPost, "/v1/roles", adminToken, access.RoleInput{Name: "Староста"}); status != http.StatusCreated {
		t.Fatalf("admin create = %d", status)
	}
	if status := call(http.MethodPost, "/v1/roles", studentToken, access.RoleInput{Name: "x"}); status != http.StatusForbidden {
		t.Fatalf("student create = %d, want 403", status)
	}
	if status := call(http.MethodPost, "/v1/roles", adminToken, map[string]any{"nme": "typo"}); status != http.StatusBadRequest {
		t.Fatalf("unknown field = %d, want 400", status)
	}
	if status := call(http.MethodGet, "/v1/me/permissions?scope=group:1", studentToken, nil); status != http.StatusOK {
		t.Fatalf("my permissions = %d", status)
	}
	if status := call(http.MethodGet, "/v1/roles", "", nil); status != http.StatusUnauthorized {
		t.Fatalf("anonymous = %d, want 401", status)
	}
}
