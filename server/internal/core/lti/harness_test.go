package lti_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/audit"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/lti"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/core/search"
	"github.com/babakapaa/lodestar/server/internal/core/structure"
	"github.com/babakapaa/lodestar/server/internal/core/tenant"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/database/dbtest"
	"github.com/babakapaa/lodestar/server/internal/platform/secretbox"
)

const (
	appURL   = "https://app.example"
	lmsIss   = "https://lms.example"
	clientID = "tool-1"
	deployID = "d1"
	lmsKeyID = "lms-key"
)

type harness struct {
	t        *testing.T
	ctx      context.Context
	tenant   string
	server   *httptest.Server
	client   *http.Client
	routes   lti.Routes
	admin    auth.Principal
	users    identity.Store
	lmsKey   *rsa.PrivateKey
	jwks     *httptest.Server
	platform lti.Platform
}

func newHarness(t *testing.T, links ...module.LinkSource) *harness {
	t.Helper()
	ctx := context.Background()
	pool := dbtest.New(t)
	sets := []database.Migrations{tenant.Migrations(), search.Migrations(), identity.Migrations(), audit.Migrations(), auth.Migrations(), access.Migrations(), structure.Migrations(), lti.Migrations()}
	if err := database.NewMigrator(pool).Apply(ctx, sets...); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	created, err := tenant.NewStore(pool).Create(ctx, "kaznu", "КазНУ")
	if err != nil {
		t.Fatalf("tenant: %v", err)
	}
	box, err := secretbox.New(base64.StdEncoding.EncodeToString(bytes.Repeat([]byte{1}, secretbox.KeySize)))
	if err != nil {
		t.Fatalf("secretbox: %v", err)
	}
	catalog, err := access.NewCatalog(lti.Permissions()...)
	if err != nil {
		t.Fatalf("catalog: %v", err)
	}
	tx := database.NewTenantTx(pool)
	acc := access.NewService(tx, catalog, structure.Tree{})
	h := &harness{t: t, ctx: ctx, tenant: created.ID, users: identity.NewStore(tx)}
	h.lmsKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("lms key: %v", err)
	}
	h.jwks = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		set := jose.JSONWebKeySet{Keys: []jose.JSONWebKey{{Key: &h.lmsKey.PublicKey, KeyID: lmsKeyID, Algorithm: string(jose.RS256), Use: "sig"}}}
		_ = json.NewEncoder(w).Encode(set)
	}))
	t.Cleanup(h.jwks.Close)

	var handler http.Handler
	h.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}))
	t.Cleanup(h.server.Close)
	authRoutes := auth.NewRoutes(pool, auth.Settings{
		Box:        box,
		Lifetimes:  auth.Lifetimes{Access: time.Minute, Refresh: time.Hour, Login: time.Minute, Grant: time.Minute},
		PublicURL:  h.server.URL,
		HTTPClient: http.DefaultClient,
	})
	h.routes = lti.NewRoutes(pool, lti.Settings{
		Box:           box,
		Access:        acc,
		Tickets:       authRoutes.Tickets,
		Authenticator: authRoutes.Authenticator,
		Links:         links,
		HTTPClient:    http.DefaultClient,
		AttemptTTL:    time.Minute,
		PublicURL:     h.server.URL,
		AppURL:        appURL,
	})
	mux := http.NewServeMux()
	authRoutes.Mount(mux)
	h.routes.Mount(mux)
	handler = mux
	h.client = &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}

	h.admin = h.user("admin", "admin@kaznu.kz")
	if _, err := access.NewBootstrap(tx).GrantTopRole(ctx, h.tenant, h.admin.UserID); err != nil {
		t.Fatalf("grant: %v", err)
	}
	h.platform, err = h.routes.Service.CreatePlatform(ctx, h.admin, lti.PlatformInput{
		Name:          "Moodle",
		Issuer:        lmsIss,
		ClientID:      clientID,
		AuthURL:       lmsIss + "/auth",
		TokenURL:      lmsIss + "/token",
		JWKSURL:       h.jwks.URL,
		DeploymentIDs: []string{deployID},
		MatchEmail:    true,
	})
	if err != nil {
		t.Fatalf("platform: %v", err)
	}
	return h
}

func (h *harness) user(subject, email string) auth.Principal {
	h.t.Helper()
	u, err := h.users.UpsertExternal(h.ctx, h.tenant, identity.ExternalIdentity{Issuer: "https://idp", Subject: subject, Email: email})
	if err != nil {
		h.t.Fatalf("user: %v", err)
	}
	return auth.Principal{TenantID: h.tenant, UserID: u.ID}
}

func (h *harness) login() url.Values {
	h.t.Helper()
	form := url.Values{"iss": {lmsIss}, "login_hint": {"hint"}, "lti_message_hint": {"msg"}, "client_id": {clientID}, "lti_deployment_id": {deployID}}
	resp, err := h.client.PostForm(h.server.URL+"/v1/lti/tenants/kaznu/login", form)
	if err != nil {
		h.t.Fatalf("login: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusFound {
		body, _ := io.ReadAll(resp.Body)
		h.t.Fatalf("login status %d: %s", resp.StatusCode, body)
	}
	target, err := url.Parse(resp.Header.Get("Location"))
	if err != nil {
		h.t.Fatalf("login location: %v", err)
	}
	return target.Query()
}

func (h *harness) token(claims map[string]any) string {
	h.t.Helper()
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: h.lmsKey},
		(&jose.SignerOptions{}).WithType("JWT").WithHeader(jose.HeaderKey("kid"), lmsKeyID))
	if err != nil {
		h.t.Fatalf("signer: %v", err)
	}
	raw, err := jwt.Signed(signer).Claims(claims).Serialize()
	if err != nil {
		h.t.Fatalf("sign: %v", err)
	}
	return raw
}

func (h *harness) claims(nonce, subject, email string, roles ...string) map[string]any {
	now := time.Now()
	return map[string]any{
		"iss":   lmsIss,
		"aud":   clientID,
		"sub":   subject,
		"exp":   now.Add(5 * time.Minute).Unix(),
		"iat":   now.Unix(),
		"nonce": nonce,
		"email": email,
		"name":  "Айгерим " + subject,
		"https://purl.imsglobal.org/spec/lti/claim/message_type":  lti.MessageResourceLink,
		"https://purl.imsglobal.org/spec/lti/claim/version":       lti.Version,
		"https://purl.imsglobal.org/spec/lti/claim/deployment_id": deployID,
		"https://purl.imsglobal.org/spec/lti/claim/roles":         roles,
		"https://purl.imsglobal.org/spec/lti/claim/context":       map[string]any{"id": "c-101", "label": "BIO-101", "title": "Основы биологии"},
		"https://purl.imsglobal.org/spec/lti-nrps/claim/namesroleservice": map[string]any{
			"context_memberships_url": lmsIss + "/memberships",
		},
		"https://purl.imsglobal.org/spec/lti-ags/claim/endpoint": map[string]any{
			"scope":     []string{"https://purl.imsglobal.org/spec/lti-ags/scope/lineitem"},
			"lineitems": lmsIss + "/lineitems",
		},
	}
}

func (h *harness) launch(state string, claims map[string]any) url.Values {
	h.t.Helper()
	resp, err := h.client.PostForm(h.server.URL+"/v1/lti/tenants/kaznu/launch", url.Values{"state": {state}, "id_token": {h.token(claims)}})
	if err != nil {
		h.t.Fatalf("launch: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusSeeOther {
		body, _ := io.ReadAll(resp.Body)
		h.t.Fatalf("launch status %d: %s", resp.StatusCode, body)
	}
	location := resp.Header.Get("Location")
	if !strings.HasPrefix(location, appURL+"/lti?") {
		h.t.Fatalf("launch location = %s", location)
	}
	target, _ := url.Parse(location)
	return target.Query()
}

func (h *harness) enter(subject, email string, roles ...string) (string, string) {
	h.t.Helper()
	q := h.login()
	result := h.launch(q.Get("state"), h.claims(q.Get("nonce"), subject, email, roles...))
	if result.Get("error") != "" {
		h.t.Fatalf("launch failed: %s", result.Get("error"))
	}
	return h.claim(result.Get("ticket")), result.Get("next")
}

func (h *harness) claim(ticket string) string {
	h.t.Helper()
	var pair auth.TokenPair
	status := h.call(http.MethodPost, "/v1/auth/token", "", map[string]string{"grantType": "launch_ticket", "ticket": ticket}, &pair)
	if status != http.StatusOK {
		h.t.Fatalf("claim status %d", status)
	}
	return pair.AccessToken
}

func (h *harness) call(method, path, bearer string, body, out any) int {
	h.t.Helper()
	var reader io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			h.t.Fatalf("encode: %v", err)
		}
		reader = bytes.NewReader(raw)
	}
	req, err := http.NewRequestWithContext(h.ctx, method, h.server.URL+path, reader)
	if err != nil {
		h.t.Fatalf("request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}
	resp, err := h.client.Do(req)
	if err != nil {
		h.t.Fatalf("do: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if out != nil && resp.StatusCode < 300 {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			h.t.Fatalf("decode: %v", err)
		}
	}
	return resp.StatusCode
}

func (h *harness) me(bearer string) identity.User {
	h.t.Helper()
	var u identity.User
	if status := h.call(http.MethodGet, "/v1/me", bearer, nil, &u); status != http.StatusOK {
		h.t.Fatalf("me status %d", status)
	}
	return u
}
