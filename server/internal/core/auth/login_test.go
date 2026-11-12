package auth_test

import (
	"net/http"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/auth/oidctest"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
)

func TestLoginEndToEnd(t *testing.T) {
	h := newHarness(t)

	pair := h.signIn()

	var me identity.User
	if status := h.do(http.MethodGet, "/v1/me", pair.AccessToken, &me); status != http.StatusOK {
		t.Fatalf("me status = %d", status)
	}
	if me.Email != "student@uni.example" || me.DisplayName != "Айгерим" {
		t.Fatalf("me = %+v", me)
	}
	var sessions []auth.Session
	if status := h.do(http.MethodGet, "/v1/sessions", pair.AccessToken, &sessions); status != http.StatusOK {
		t.Fatalf("sessions status = %d", status)
	}
	if len(sessions) != 1 || !sessions[0].Current || sessions[0].DeviceName != "Pixel 9" {
		t.Fatalf("sessions = %+v", sessions)
	}
}

func TestSameAccountReusesUser(t *testing.T) {
	h := newHarness(t)
	first := h.signIn()
	h.idp.SignIn(oidctest.User{Subject: "student-1", Email: "new@uni.example", Name: "Айгерим С."})
	second := h.signIn()

	var a, b identity.User
	h.do(http.MethodGet, "/v1/me", first.AccessToken, &a)
	h.do(http.MethodGet, "/v1/me", second.AccessToken, &b)
	if a.ID != b.ID || b.Email != "new@uni.example" {
		t.Fatalf("first %+v, second %+v", a, b)
	}
}

func TestLoginRejectsUnknownRedirect(t *testing.T) {
	h := newHarness(t)
	_, challenge := newVerifier()

	resp, err := h.client.Get(h.server.URL + "/v1/auth/kaznu/campus/login?code_challenge_method=S256&redirect_uri=https://evil.example&code_challenge=" + challenge)
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", resp.StatusCode)
	}
}

func TestExchangeRequiresMatchingVerifier(t *testing.T) {
	h := newHarness(t)
	_, challenge := newVerifier()
	back := h.login(challenge)
	wrong, _ := newVerifier()

	status := h.token(auth.TokenRequest{
		GrantType:    "authorization_code",
		Code:         back.Get("code"),
		CodeVerifier: wrong,
		RedirectURI:  appRedirect,
	}, nil)
	if status != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", status)
	}
}

func TestProvidersAreListed(t *testing.T) {
	h := newHarness(t)

	var list []auth.ProviderInfo
	resp, err := h.client.Get(h.server.URL + "/v1/auth/kaznu/providers")
	if err != nil {
		t.Fatalf("providers: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if status := decode(t, resp, &list); status != http.StatusOK {
		t.Fatalf("status = %d", status)
	}
	if len(list) != 1 || list[0].Key != "campus" {
		t.Fatalf("providers = %+v", list)
	}
}

func TestUnknownTenant(t *testing.T) {
	h := newHarness(t)

	resp, err := h.client.Get(h.server.URL + "/v1/auth/missing/providers")
	if err != nil {
		t.Fatalf("providers: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if status := decode(t, resp, nil); status != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", status)
	}
}
