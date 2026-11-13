package lti

import (
	"fmt"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"

	"github.com/babakapaa/lodestar/server/internal/core/module"
)

func deepLinkResponse(p Platform, d DeepLink, targets []module.LinkTarget, key Key) (string, error) {
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: key.Private},
		(&jose.SignerOptions{}).WithType("JWT").WithHeader(jose.HeaderKey("kid"), key.ID))
	if err != nil {
		return "", fmt.Errorf("lti signer: %w", err)
	}
	nonce, err := randomString()
	if err != nil {
		return "", err
	}
	items := make([]map[string]any, len(targets))
	for i, t := range targets {
		items[i] = map[string]any{
			"type":   "ltiResourceLink",
			"title":  t.Title,
			"custom": map[string]string{customKind: t.Kind, customID: t.ID},
		}
	}
	extra := map[string]any{
		"nonce": nonce,
		"https://purl.imsglobal.org/spec/lti/claim/message_type":     messageDeepAnswer,
		"https://purl.imsglobal.org/spec/lti/claim/version":          Version,
		"https://purl.imsglobal.org/spec/lti/claim/deployment_id":    d.DeploymentID,
		"https://purl.imsglobal.org/spec/lti-dl/claim/content_items": items,
	}
	if d.Data != "" {
		extra["https://purl.imsglobal.org/spec/lti-dl/claim/data"] = d.Data
	}
	now := time.Now()
	claims := jwt.Claims{
		Issuer:   p.ClientID,
		Audience: jwt.Audience{p.Issuer},
		IssuedAt: jwt.NewNumericDate(now),
		Expiry:   jwt.NewNumericDate(now.Add(assertionLife)),
	}
	raw, err := jwt.Signed(signer).Claims(claims).Claims(extra).Serialize()
	if err != nil {
		return "", fmt.Errorf("lti deep linking response: %w", err)
	}
	return raw, nil
}
