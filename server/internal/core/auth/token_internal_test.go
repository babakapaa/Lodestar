package auth

import (
	"errors"
	"testing"
)

func TestTokenCarriesTenant(t *testing.T) {
	tenantID := "0b6f3c1e-8a2d-4f5b-9c7e-1d2a3b4c5d6e"
	token, err := newToken(accessPrefix, tenantID)
	if err != nil {
		t.Fatalf("newToken: %v", err)
	}
	parsed, err := parseToken(accessPrefix, token.Value)
	if err != nil {
		t.Fatalf("parseToken: %v", err)
	}
	if parsed.TenantID != tenantID || string(parsed.Hash) != string(token.Hash) {
		t.Fatalf("parsed %+v, want tenant %s", parsed, tenantID)
	}
}

func TestParseTokenRejectsWrongKind(t *testing.T) {
	token, err := newToken(refreshPrefix, "0b6f3c1e-8a2d-4f5b-9c7e-1d2a3b4c5d6e")
	if err != nil {
		t.Fatalf("newToken: %v", err)
	}
	if _, err := parseToken(accessPrefix, token.Value); !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("err = %v, want ErrInvalidToken", err)
	}
}
