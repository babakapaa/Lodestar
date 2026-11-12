package demoidp_test

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/platform/demoidp"
)

const (
	verifier = "a-long-enough-verifier-for-pkce-tests-0123456789"
	callback = "http://lodestar.test/v1/auth/demo/demo/callback"
)

func start(t *testing.T) *httptest.Server {
	t.Helper()
	p, err := demoidp.New(demoidp.Config{
		ClientID:     "lodestar-demo",
		ClientSecret: "secret",
		Redirects:    []string{callback},
		Users: []demoidp.User{
			{Subject: "s1", Email: "student@demo.invalid", Name: "Айгерим", Title: "Студентка"},
			{Subject: "t1", Email: "teacher@demo.invalid", Name: "Ерлан", Title: "Преподаватель"},
		},
	})
	if err != nil {
		t.Fatalf("new: %v", err)
	}
	srv := httptest.NewServer(p.Handler())
	t.Cleanup(srv.Close)
	p.SetIssuer(srv.URL)
	return srv
}

func authorizeQuery(redirect string) url.Values {
	sum := sha256.Sum256([]byte(verifier))
	return url.Values{
		"client_id":             {"lodestar-demo"},
		"redirect_uri":          {redirect},
		"response_type":         {"code"},
		"state":                 {"st"},
		"nonce":                 {"n"},
		"code_challenge":        {base64.RawURLEncoding.EncodeToString(sum[:])},
		"code_challenge_method": {"S256"},
	}
}

func noRedirects() *http.Client {
	return &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
}

func TestChoosingAPersonSignsThemIn(t *testing.T) {
	srv := start(t)
	q := authorizeQuery(callback)

	page, err := http.Get(srv.URL + "/authorize?" + q.Encode())
	if err != nil {
		t.Fatalf("page: %v", err)
	}
	body, _ := io.ReadAll(page.Body)
	_ = page.Body.Close()
	if !strings.Contains(string(body), "Ерлан") || !strings.Contains(string(body), "Преподаватель") {
		t.Fatalf("page does not list people: %s", body)
	}

	form := authorizeQuery(callback)
	form.Set("subject", "t1")
	resp, err := noRedirects().PostForm(srv.URL+"/authorize", form)
	if err != nil {
		t.Fatalf("choose: %v", err)
	}
	_ = resp.Body.Close()
	back, err := url.Parse(resp.Header.Get("Location"))
	if err != nil || resp.StatusCode != http.StatusFound || back.Query().Get("state") != "st" {
		t.Fatalf("redirect = %d %q", resp.StatusCode, resp.Header.Get("Location"))
	}

	token, err := http.PostForm(srv.URL+"/token", url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {back.Query().Get("code")},
		"redirect_uri":  {callback},
		"code_verifier": {verifier},
		"client_id":     {"lodestar-demo"},
		"client_secret": {"secret"},
	})
	if err != nil {
		t.Fatalf("token: %v", err)
	}
	defer func() { _ = token.Body.Close() }()
	var out struct {
		IDToken string `json:"id_token"`
	}
	if err := json.NewDecoder(token.Body).Decode(&out); err != nil || token.StatusCode != http.StatusOK {
		t.Fatalf("token status %d: %v", token.StatusCode, err)
	}
	parts := strings.Split(out.IDToken, ".")
	claims, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil || !strings.Contains(string(claims), `"sub":"t1"`) || !strings.Contains(string(claims), "teacher@demo.invalid") {
		t.Fatalf("claims = %s (%v)", claims, err)
	}
}

func TestUnknownRedirectIsRefused(t *testing.T) {
	srv := start(t)

	resp, err := http.Get(srv.URL + "/authorize?" + authorizeQuery("https://evil.example/cb").Encode())
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", resp.StatusCode)
	}
}

func TestCodeWorksOnce(t *testing.T) {
	srv := start(t)
	form := authorizeQuery(callback)
	form.Set("subject", "s1")
	resp, err := noRedirects().PostForm(srv.URL+"/authorize", form)
	if err != nil {
		t.Fatalf("choose: %v", err)
	}
	_ = resp.Body.Close()
	back, _ := url.Parse(resp.Header.Get("Location"))
	exchange := url.Values{
		"code":          {back.Query().Get("code")},
		"redirect_uri":  {callback},
		"code_verifier": {verifier},
		"client_id":     {"lodestar-demo"},
		"client_secret": {"secret"},
	}
	for i, want := range []int{http.StatusOK, http.StatusBadRequest} {
		r, err := http.PostForm(srv.URL+"/token", exchange)
		if err != nil {
			t.Fatalf("token %d: %v", i, err)
		}
		_ = r.Body.Close()
		if r.StatusCode != want {
			t.Errorf("exchange %d status = %d, want %d", i, r.StatusCode, want)
		}
	}
}
