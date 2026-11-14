package tester_test

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
	"strings"
	"testing"
	"time"

	"github.com/babakapaa/lodestar/server/internal/app"
	"github.com/babakapaa/lodestar/server/internal/app/demo"
	"github.com/babakapaa/lodestar/server/internal/app/tester"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/platform/config"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/database/dbtest"
	"github.com/babakapaa/lodestar/server/internal/platform/secretbox"
)

const appRedirect = "http://localhost:43823/auth"

type world struct {
	t      *testing.T
	url    string
	client *http.Client
}

func start(t *testing.T) world {
	t.Helper()
	ctx := context.Background()
	pool := dbtest.New(t)
	registry, err := module.NewRegistry(app.Modules()...)
	if err != nil {
		t.Fatalf("registry: %v", err)
	}
	if err := database.NewMigrator(pool).Apply(ctx, append(app.CoreMigrations(), registry.Migrations()...)...); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	box, err := secretbox.New(base64.StdEncoding.EncodeToString(bytes.Repeat([]byte{7}, secretbox.KeySize)))
	if err != nil {
		t.Fatalf("box: %v", err)
	}
	srv := httptest.NewUnstartedServer(nil)
	publicURL := "http://" + srv.Listener.Addr().String()
	cfg := config.Config{
		AccessTTL:    time.Minute,
		RefreshTTL:   time.Hour,
		PublicURL:    publicURL,
		AppURL:       publicURL,
		AppRedirects: []string{appRedirect},
		Tester:       true,
	}
	if _, err := tester.Start(ctx, cfg, pool, box, registry); err != nil {
		t.Fatalf("start: %v", err)
	}
	mode, err := tester.Start(ctx, cfg, pool, box, registry)
	if err != nil {
		t.Fatalf("restart: %v", err)
	}
	server, err := app.NewServer("test", cfg, box, pool, registry, app.Integrations{LoginClient: mode.LoginClient()})
	if err != nil {
		t.Fatalf("server: %v", err)
	}
	srv.Config.Handler = mode.Wrap(server.Handler)
	srv.Start()
	t.Cleanup(srv.Close)
	client := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	return world{t: t, url: publicURL, client: client}
}

type reply struct {
	status   int
	body     []byte
	location string
}

func (r reply) decode(t *testing.T, v any) {
	t.Helper()
	if err := json.Unmarshal(r.body, v); err != nil {
		t.Fatalf("decode %d %s: %v", r.status, r.body, err)
	}
}

func (w world) do(req *http.Request) reply {
	w.t.Helper()
	resp, err := w.client.Do(req)
	if err != nil {
		w.t.Fatalf("%s %s: %v", req.Method, req.URL, err)
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		w.t.Fatalf("read %s: %v", req.URL, err)
	}
	return reply{status: resp.StatusCode, body: body, location: resp.Header.Get("Location")}
}

func (w world) post(path string, body any) reply {
	w.t.Helper()
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, w.url+path, bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	return w.do(req)
}

func (w world) get(target, token string) reply {
	w.t.Helper()
	req, _ := http.NewRequest(http.MethodGet, target, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return w.do(req)
}

func (w world) me(token string) string {
	w.t.Helper()
	var me struct {
		Email string `json:"email"`
	}
	w.get(w.url+"/v1/me", token).decode(w.t, &me)
	return me.Email
}

func (w world) token(r reply) string {
	w.t.Helper()
	var pair struct {
		AccessToken string `json:"accessToken"`
	}
	if r.status != http.StatusOK {
		w.t.Fatalf("token = %d %s", r.status, r.body)
	}
	r.decode(w.t, &pair)
	return pair.AccessToken
}

func TestSignInAsAnyDemoPerson(t *testing.T) {
	w := start(t)

	var info tester.Info
	w.get(w.url+"/v1/tester", "").decode(t, &info)
	if info.Tenant != demo.TenantSlug || len(info.People) != len(demo.People()) {
		t.Fatalf("info = %+v", info)
	}
	for _, person := range info.People {
		token := w.token(w.post("/v1/tester/sign-in", tester.SignInRequest{Person: person.Key}))
		if got := w.me(token); got != person.Email {
			t.Errorf("signed in as %s, want %s", got, person.Email)
		}
	}
	if r := w.post("/v1/tester/sign-in", tester.SignInRequest{Person: "stranger"}); r.status != http.StatusNotFound {
		t.Errorf("unknown person = %d", r.status)
	}
}

func TestBuiltInDemoLoginGoesThroughOIDC(t *testing.T) {
	w := start(t)
	verifier := rand.Text() + rand.Text()
	sum := sha256.Sum256([]byte(verifier))
	login := w.url + "/v1/auth/demo/demo/login?" + url.Values{
		"redirect_uri":          {appRedirect},
		"code_challenge":        {base64.RawURLEncoding.EncodeToString(sum[:])},
		"code_challenge_method": {"S256"},
	}.Encode()

	toIdP := w.get(login, "")
	authorize, err := url.Parse(toIdP.location)
	if toIdP.status != http.StatusFound || err != nil || !strings.HasPrefix(toIdP.location, w.url+tester.IdPPath+"/authorize") {
		t.Fatalf("login = %d %q", toIdP.status, toIdP.location)
	}
	if page := w.get(toIdP.location, ""); page.status != http.StatusOK || !strings.Contains(string(page.body), demo.Teacher.Name) {
		t.Fatalf("choice page = %d", page.status)
	}
	form := authorize.Query()
	form.Set("subject", demo.Teacher.Subject)
	chosen := w.do(mustForm(t, w.url+tester.IdPPath+"/authorize", form))
	callback := w.get(chosen.location, "")
	back, err := url.Parse(callback.location)
	if callback.status != http.StatusFound || err != nil || back.Query().Get("code") == "" {
		t.Fatalf("callback = %d %q", callback.status, callback.location)
	}
	token := w.token(w.post("/v1/auth/token", map[string]string{
		"grantType":    "authorization_code",
		"code":         back.Query().Get("code"),
		"codeVerifier": verifier,
		"redirectUri":  appRedirect,
	}))
	if got := w.me(token); got != demo.Teacher.Email {
		t.Errorf("signed in as %s", got)
	}
}

func mustForm(t *testing.T, target string, form url.Values) *http.Request {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, target, strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}
