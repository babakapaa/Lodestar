package auth_test

import (
	"net/http"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
)

func TestRefreshRotatesTokens(t *testing.T) {
	h := newHarness(t)
	first := h.signIn()

	var second auth.TokenPair
	status := h.token(auth.TokenRequest{GrantType: "refresh_token", RefreshToken: first.RefreshToken}, &second)
	if status != http.StatusOK {
		t.Fatalf("refresh status = %d", status)
	}
	if second.RefreshToken == first.RefreshToken {
		t.Fatal("refresh token was not rotated")
	}
	if status := h.do(http.MethodGet, "/v1/me", first.AccessToken, nil); status != http.StatusUnauthorized {
		t.Fatalf("old access status = %d, want 401", status)
	}
	if status := h.do(http.MethodGet, "/v1/me", second.AccessToken, nil); status != http.StatusOK {
		t.Fatalf("new access status = %d, want 200", status)
	}
}

func TestReusedRefreshRevokesSession(t *testing.T) {
	h := newHarness(t)
	first := h.signIn()
	var second auth.TokenPair
	h.token(auth.TokenRequest{GrantType: "refresh_token", RefreshToken: first.RefreshToken}, &second)

	status := h.token(auth.TokenRequest{GrantType: "refresh_token", RefreshToken: first.RefreshToken}, nil)
	if status != http.StatusBadRequest {
		t.Fatalf("reuse status = %d, want 400", status)
	}
	if status := h.do(http.MethodGet, "/v1/me", second.AccessToken, nil); status != http.StatusUnauthorized {
		t.Fatalf("session survived refresh reuse, status = %d", status)
	}
}

func TestLogoutEndsSession(t *testing.T) {
	h := newHarness(t)
	pair := h.signIn()

	if status := h.do(http.MethodPost, "/v1/auth/logout", pair.AccessToken, nil); status != http.StatusNoContent {
		t.Fatalf("logout status = %d", status)
	}
	if status := h.do(http.MethodGet, "/v1/me", pair.AccessToken, nil); status != http.StatusUnauthorized {
		t.Fatalf("me after logout = %d, want 401", status)
	}
}

func TestRevokeOtherDevice(t *testing.T) {
	h := newHarness(t)
	phone := h.signIn()
	laptop := h.signIn()
	var sessions []auth.Session
	h.do(http.MethodGet, "/v1/sessions", laptop.AccessToken, &sessions)
	var phoneID string
	for _, s := range sessions {
		if !s.Current {
			phoneID = s.ID
		}
	}

	if status := h.do(http.MethodDelete, "/v1/sessions/"+phoneID, laptop.AccessToken, nil); status != http.StatusNoContent {
		t.Fatalf("revoke status = %d", status)
	}
	if status := h.do(http.MethodGet, "/v1/me", phone.AccessToken, nil); status != http.StatusUnauthorized {
		t.Fatalf("phone still signed in, status = %d", status)
	}
	if status := h.do(http.MethodDelete, "/v1/sessions/"+phoneID, laptop.AccessToken, nil); status != http.StatusNotFound {
		t.Fatalf("second revoke status = %d, want 404", status)
	}
}

func TestMissingTokenIsUnauthorized(t *testing.T) {
	h := newHarness(t)

	if status := h.do(http.MethodGet, "/v1/me", "lsa_garbage", nil); status != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", status)
	}
}
