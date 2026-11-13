package sandbox

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"strings"

	"github.com/google/uuid"
)

const tokenPrefix = "lsn_"

type NodeToken struct {
	TenantID string
	raw      string
}

func NewNodeToken(tenantID string) (NodeToken, error) {
	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		return NodeToken{}, err
	}
	return NodeToken{TenantID: tenantID, raw: tokenPrefix + tenantID + "_" + base64.RawURLEncoding.EncodeToString(secret)}, nil
}

func ParseNodeToken(raw string) (NodeToken, error) {
	rest, ok := strings.CutPrefix(raw, tokenPrefix)
	if !ok {
		return NodeToken{}, ErrBadNodeToken
	}
	tenantID, secret, ok := strings.Cut(rest, "_")
	if !ok || secret == "" {
		return NodeToken{}, ErrBadNodeToken
	}
	if _, err := uuid.Parse(tenantID); err != nil {
		return NodeToken{}, ErrBadNodeToken
	}
	return NodeToken{TenantID: tenantID, raw: raw}, nil
}

func (t NodeToken) String() string {
	return t.raw
}

func (t NodeToken) Hash() []byte {
	sum := sha256.Sum256([]byte(t.raw))
	return sum[:]
}
