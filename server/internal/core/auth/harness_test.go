package auth_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/auth/oidctest"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/tenant"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/database/dbtest"
	"github.com/babakapaa/lodestar/server/internal/platform/secretbox"
)

const (
	appRedirect = "lodestar-test:/auth"
	appScheme   = "lodestar-test"
)

type harness struct {
	t      *testing.T
	idp    *oidctest.Server
	server *httptest.Server
	client *http.Client
}

func newHarness(t *testing.T) *harness {
	t.Helper()
	ctx := context.Background()
	pool := dbtest.New(t)
	err := database.NewMigrator(pool).Apply(ctx, tenant.Migrations(), identity.Migrations(), auth.Migrations())
	if err != nil {
		t.Fatalf("migrate: %v", err)
	}
	created, err := tenant.NewStore(pool).Create(ctx, "kaznu", "КазНУ")
	if err != nil {
		t.Fatalf("create tenant: %v", err)
	}
	box, err := secretbox.New(base64.StdEncoding.EncodeToString(bytes.Repeat([]byte{1}, secretbox.KeySize)))
	if err != nil {
		t.Fatalf("secretbox: %v", err)
	}
	idp := oidctest.New(t)
	_, err = auth.NewProviderStore(database.NewTenantTx(pool), box).Put(ctx, auth.Provider{
		TenantID:     created.ID,
		Key:          "campus",
		DisplayName:  "Кампус",
		Issuer:       idp.URL,
		ClientID:     idp.ClientID,
		ClientSecret: idp.ClientSecret,
	})
	if err != nil {
		t.Fatalf("put provider: %v", err)
	}

	var handler http.Handler
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}))
	t.Cleanup(server.Close)
	mux := http.NewServeMux()
	auth.NewRoutes(pool, auth.Settings{
		Box:        box,
		Lifetimes:  auth.Lifetimes{Access: time.Minute, Refresh: time.Hour, Login: time.Minute, Grant: time.Minute},
		Redirects:  []string{appRedirect},
		PublicURL:  server.URL,
		HTTPClient: http.DefaultClient,
	}).Mount(mux)
	handler = mux

	client := &http.Client{CheckRedirect: func(req *http.Request, _ []*http.Request) error {
		if req.URL.Scheme == appScheme {
			return http.ErrUseLastResponse
		}
		return nil
	}}
	return &harness{t: t, idp: idp, server: server, client: client}
}

func newVerifier() (string, string) {
	verifier := rand.Text() + rand.Text()
	sum := sha256.Sum256([]byte(verifier))
	return verifier, base64.RawURLEncoding.EncodeToString(sum[:])
}

func (h *harness) login(challenge string) url.Values {
	h.t.Helper()
	q := url.Values{
		"redirect_uri":          {appRedirect},
		"code_challenge":        {challenge},
		"code_challenge_method": {"S256"},
		"device":                {"Pixel 9"},
	}
	resp, err := h.client.Get(h.server.URL + "/v1/auth/kaznu/campus/login?" + q.Encode())
	if err != nil {
		h.t.Fatalf("login: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusFound {
		body, _ := io.ReadAll(resp.Body)
		h.t.Fatalf("login status = %d: %s", resp.StatusCode, body)
	}
	target, err := url.Parse(resp.Header.Get("Location"))
	if err != nil {
		h.t.Fatalf("parse app redirect: %v", err)
	}
	return target.Query()
}

func (h *harness) signIn() auth.TokenPair {
	h.t.Helper()
	verifier, challenge := newVerifier()
	back := h.login(challenge)
	var pair auth.TokenPair
	status := h.token(auth.TokenRequest{
		GrantType:    "authorization_code",
		Code:         back.Get("code"),
		CodeVerifier: verifier,
		RedirectURI:  appRedirect,
	}, &pair)
	if status != http.StatusOK {
		h.t.Fatalf("exchange status = %d", status)
	}
	return pair
}

func (h *harness) token(body auth.TokenRequest, out any) int {
	h.t.Helper()
	data, err := json.Marshal(body)
	if err != nil {
		h.t.Fatalf("marshal: %v", err)
	}
	resp, err := h.client.Post(h.server.URL+"/v1/auth/token", "application/json", bytes.NewReader(data))
	if err != nil {
		h.t.Fatalf("token: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	return decode(h.t, resp, out)
}

func (h *harness) do(method, path, access string, out any) int {
	h.t.Helper()
	req, err := http.NewRequest(method, h.server.URL+path, nil)
	if err != nil {
		h.t.Fatalf("request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+access)
	resp, err := h.client.Do(req)
	if err != nil {
		h.t.Fatalf("%s %s: %v", method, path, err)
	}
	defer func() { _ = resp.Body.Close() }()
	return decode(h.t, resp, out)
}

func decode(t *testing.T, resp *http.Response, out any) int {
	t.Helper()
	if out != nil && resp.StatusCode < 300 {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			t.Fatalf("decode: %v", err)
		}
	}
	return resp.StatusCode
}
