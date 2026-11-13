package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

const (
	accessPrefix  = "lsa_"
	refreshPrefix = "lsr_"
	grantPrefix   = "lsg_"
	ticketPrefix  = "lst_"
	tenantBytes   = 16
	secretBytes   = 32
)

type Token struct {
	Value    string
	TenantID string
	Hash     []byte
}

func newToken(prefix, tenantID string) (Token, error) {
	tenant, err := uuidBytes(tenantID)
	if err != nil {
		return Token{}, err
	}
	secret := make([]byte, secretBytes)
	if _, err := rand.Read(secret); err != nil {
		return Token{}, fmt.Errorf("token: %w", err)
	}
	value := prefix + base64.RawURLEncoding.EncodeToString(append(tenant, secret...))
	return Token{Value: value, TenantID: tenantID, Hash: hashToken(value)}, nil
}

func parseToken(prefix, value string) (Token, error) {
	encoded, ok := strings.CutPrefix(value, prefix)
	if !ok {
		return Token{}, ErrInvalidToken
	}
	raw, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil || len(raw) != tenantBytes+secretBytes {
		return Token{}, ErrInvalidToken
	}
	return Token{Value: value, TenantID: uuidString(raw[:tenantBytes]), Hash: hashToken(value)}, nil
}

func hashToken(value string) []byte {
	sum := sha256.Sum256([]byte(value))
	return sum[:]
}

func uuidBytes(id string) ([]byte, error) {
	raw, err := hex.DecodeString(strings.ReplaceAll(id, "-", ""))
	if err != nil || len(raw) != tenantBytes {
		return nil, fmt.Errorf("tenant id %q is not a uuid", id)
	}
	return raw, nil
}

func uuidString(raw []byte) string {
	h := hex.EncodeToString(raw)
	return h[0:8] + "-" + h[8:12] + "-" + h[12:16] + "-" + h[16:20] + "-" + h[20:32]
}
