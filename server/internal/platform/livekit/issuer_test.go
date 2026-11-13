package livekit_test

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"

	"github.com/babakapaa/lodestar/server/internal/platform/livekit"
)

const secret = "0123456789abcdef0123456789abcdef"

func TestTokenCarriesRoomAndSources(t *testing.T) {
	now := time.Date(2026, 11, 13, 9, 0, 0, 0, time.UTC)
	issuer := livekit.NewIssuer(livekit.Config{URL: "wss://calls", APIKey: "key", APISecret: secret})
	raw, err := issuer.Token("u1", "Айгерим", livekit.Grant{Room: "t_c", ScreenShare: false}, time.Hour, now)
	if err != nil {
		t.Fatalf("token: %v", err)
	}
	parsed, err := jwt.ParseSigned(raw, []jose.SignatureAlgorithm{jose.HS256})
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	var std jwt.Claims
	var extra struct {
		Name  string `json:"name"`
		Video struct {
			Room    string   `json:"room"`
			Join    bool     `json:"roomJoin"`
			Sources []string `json:"canPublishSources"`
		} `json:"video"`
	}
	if err := parsed.Claims([]byte(secret), &std, &extra); err != nil {
		t.Fatalf("claims: %v", err)
	}
	if std.Issuer != "key" || std.Subject != "u1" || extra.Name != "Айгерим" || extra.Video.Room != "t_c" || !extra.Video.Join {
		t.Fatalf("claims = %+v %+v", std, extra)
	}
	if slices.Contains(extra.Video.Sources, "screen_share") {
		t.Fatalf("sources = %v", extra.Video.Sources)
	}
}

func TestWebhookSignature(t *testing.T) {
	now := time.Now()
	issuer := livekit.NewIssuer(livekit.Config{URL: "wss://calls", APIKey: "key", APISecret: secret})
	body := []byte(`{"event":"room_finished"}`)
	sum := sha256.Sum256(body)
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS256, Key: []byte(secret)}, nil)
	if err != nil {
		t.Fatalf("signer: %v", err)
	}
	auth, err := jwt.Signed(signer).Claims(jwt.Claims{Issuer: "key", Expiry: jwt.NewNumericDate(now.Add(time.Minute))}).
		Claims(map[string]any{"sha256": base64.StdEncoding.EncodeToString(sum[:])}).Serialize()
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	if err := issuer.Verify(auth, body, now); err != nil {
		t.Fatalf("verify: %v", err)
	}
	if err := issuer.Verify(auth, []byte(`{}`), now); !errors.Is(err, livekit.ErrBadWebhook) {
		t.Fatalf("tampered: %v", err)
	}
}
