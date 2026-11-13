package livekit

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
)

type Issuer struct {
	cfg Config
}

func NewIssuer(cfg Config) *Issuer {
	return &Issuer{cfg: cfg}
}

func (i *Issuer) URL() string {
	return i.cfg.URL
}

func (i *Issuer) Token(identity, name string, g Grant, ttl time.Duration, now time.Time) (string, error) {
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS256, Key: []byte(i.cfg.APISecret)}, (&jose.SignerOptions{}).WithType("JWT"))
	if err != nil {
		return "", fmt.Errorf("livekit signer: %w", err)
	}
	sources := []string{"camera", "microphone"}
	if g.ScreenShare {
		sources = append(sources, "screen_share", "screen_share_audio")
	}
	std := jwt.Claims{
		Issuer:    i.cfg.APIKey,
		Subject:   identity,
		NotBefore: jwt.NewNumericDate(now.Add(-time.Minute)),
		Expiry:    jwt.NewNumericDate(now.Add(ttl)),
	}
	extra := map[string]any{
		"name": name,
		"video": videoClaims{
			Room: g.Room, RoomJoin: true, CanPublish: true, CanSubscribe: true, CanPublishData: true, CanPublishSources: sources,
		},
	}
	token, err := jwt.Signed(signer).Claims(std).Claims(extra).Serialize()
	if err != nil {
		return "", fmt.Errorf("livekit token: %w", err)
	}
	return token, nil
}

func (i *Issuer) Verify(authorization string, body []byte, now time.Time) error {
	parsed, err := jwt.ParseSigned(authorization, []jose.SignatureAlgorithm{jose.HS256})
	if err != nil {
		return ErrBadWebhook
	}
	var std jwt.Claims
	var extra struct {
		SHA256 string `json:"sha256"`
	}
	if err := parsed.Claims([]byte(i.cfg.APISecret), &std, &extra); err != nil {
		return ErrBadWebhook
	}
	if err := std.ValidateWithLeeway(jwt.Expected{Issuer: i.cfg.APIKey, Time: now}, time.Minute); err != nil {
		return ErrBadWebhook
	}
	sum := sha256.Sum256(body)
	if extra.SHA256 != base64.StdEncoding.EncodeToString(sum[:]) {
		return ErrBadWebhook
	}
	return nil
}
