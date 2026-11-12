package auth

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"

	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/tenant"
)

const maxDeviceName = 100

type LoginService struct {
	tenants   tenant.Store
	users     identity.Store
	providers ProviderStore
	attempts  AttemptStore
	grants    GrantStore
	sessions  SessionStore
	discovery *Discovery
	redirects Redirects
	publicURL string
}

type LoginDeps struct {
	Tenants   tenant.Store
	Users     identity.Store
	Providers ProviderStore
	Attempts  AttemptStore
	Grants    GrantStore
	Sessions  SessionStore
	Discovery *Discovery
	Redirects Redirects
	PublicURL string
}

func NewLoginService(d LoginDeps) LoginService {
	return LoginService{
		tenants:   d.Tenants,
		users:     d.Users,
		providers: d.Providers,
		attempts:  d.Attempts,
		grants:    d.Grants,
		sessions:  d.Sessions,
		discovery: d.Discovery,
		redirects: d.Redirects,
		publicURL: d.PublicURL,
	}
}

func (s LoginService) Providers(ctx context.Context, tenantSlug string) ([]ProviderInfo, error) {
	t, err := s.tenants.BySlug(ctx, tenantSlug)
	if err != nil {
		return nil, err
	}
	return s.providers.List(ctx, t.ID)
}

func (s LoginService) Begin(ctx context.Context, req LoginRequest) (string, error) {
	if !s.redirects.Allows(req.RedirectURI) {
		return "", ErrRedirectNotAllowed
	}
	if !validChallenge(req.CodeChallenge) {
		return "", ErrInvalidChallenge
	}
	t, err := s.tenants.BySlug(ctx, req.TenantSlug)
	if err != nil {
		return "", err
	}
	p, err := s.providers.ByKey(ctx, t.ID, req.ProviderKey)
	if err != nil {
		return "", err
	}
	oauth, _, err := s.client(ctx, p, req.TenantSlug)
	if err != nil {
		return "", err
	}
	attempt := LoginAttempt{
		ProviderID:    p.ID,
		Verifier:      oauth2.GenerateVerifier(),
		RedirectURI:   req.RedirectURI,
		CodeChallenge: req.CodeChallenge,
		DeviceName:    deviceName(req.DeviceName),
	}
	if attempt.State, err = randomString(); err != nil {
		return "", err
	}
	if attempt.Nonce, err = randomString(); err != nil {
		return "", err
	}
	if err := s.attempts.Create(ctx, t.ID, attempt); err != nil {
		return "", err
	}
	return oauth.AuthCodeURL(attempt.State, oidc.Nonce(attempt.Nonce), oauth2.S256ChallengeOption(attempt.Verifier)), nil
}

func (s LoginService) Complete(ctx context.Context, cb Callback) (string, error) {
	t, err := s.tenants.BySlug(ctx, cb.TenantSlug)
	if err != nil {
		return "", err
	}
	attempt, err := s.attempts.Take(ctx, t.ID, cb.State)
	if err != nil {
		return "", err
	}
	user, err := s.identify(ctx, t.ID, cb, attempt)
	if err != nil {
		return appRedirect(attempt.RedirectURI, "error", "login_failed"), fmt.Errorf("%w: %w", ErrLoginFailed, err)
	}
	code, err := s.grants.Create(ctx, LoginGrant{
		TenantID:      t.ID,
		UserID:        user.ID,
		RedirectURI:   attempt.RedirectURI,
		CodeChallenge: attempt.CodeChallenge,
		DeviceName:    attempt.DeviceName,
	})
	if err != nil {
		return "", err
	}
	return appRedirect(attempt.RedirectURI, "code", code), nil
}

func (s LoginService) Exchange(ctx context.Context, code, verifier, redirectURI string) (TokenPair, error) {
	g, err := s.grants.Take(ctx, code)
	if err != nil {
		return TokenPair{}, err
	}
	if g.RedirectURI != redirectURI || !challengeMatches(g.CodeChallenge, verifier) {
		return TokenPair{}, ErrInvalidToken
	}
	return s.sessions.Start(ctx, g.TenantID, g.UserID, g.DeviceName)
}

func (s LoginService) identify(ctx context.Context, tenantID string, cb Callback, attempt LoginAttempt) (identity.User, error) {
	if cb.Error != "" {
		return identity.User{}, fmt.Errorf("provider returned %s", cb.Error)
	}
	p, err := s.providers.ByKey(ctx, tenantID, cb.ProviderKey)
	if err != nil {
		return identity.User{}, err
	}
	if p.ID != attempt.ProviderID {
		return identity.User{}, ErrInvalidState
	}
	oauth, verifier, err := s.client(ctx, p, cb.TenantSlug)
	if err != nil {
		return identity.User{}, err
	}
	httpCtx := s.discovery.Context(ctx)
	token, err := oauth.Exchange(httpCtx, cb.Code, oauth2.VerifierOption(attempt.Verifier))
	if err != nil {
		return identity.User{}, fmt.Errorf("exchange code: %w", err)
	}
	raw, ok := token.Extra("id_token").(string)
	if !ok {
		return identity.User{}, fmt.Errorf("provider returned no id_token")
	}
	idToken, err := verifier.Verify(httpCtx, raw)
	if err != nil {
		return identity.User{}, fmt.Errorf("verify id_token: %w", err)
	}
	if idToken.Nonce != attempt.Nonce {
		return identity.User{}, fmt.Errorf("id_token nonce mismatch")
	}
	var claims idClaims
	if err := idToken.Claims(&claims); err != nil {
		return identity.User{}, fmt.Errorf("read claims: %w", err)
	}
	return s.users.UpsertExternal(ctx, tenantID, identity.ExternalIdentity{
		Issuer:      idToken.Issuer,
		Subject:     idToken.Subject,
		Email:       claims.Email,
		DisplayName: claims.Name,
	})
}

func (s LoginService) client(ctx context.Context, p Provider, tenantSlug string) (oauth2.Config, *oidc.IDTokenVerifier, error) {
	provider, err := s.discovery.Provider(ctx, p.Issuer)
	if err != nil {
		return oauth2.Config{}, nil, err
	}
	oauth := oauth2.Config{
		ClientID:     p.ClientID,
		ClientSecret: p.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  CallbackURL(s.publicURL, tenantSlug, p.Key),
		Scopes:       p.Scopes,
	}
	return oauth, provider.Verifier(&oidc.Config{ClientID: p.ClientID}), nil
}

func appRedirect(uri, key, value string) string {
	target, err := url.Parse(uri)
	if err != nil {
		return uri
	}
	q := target.Query()
	q.Set(key, value)
	target.RawQuery = q.Encode()
	return target.String()
}

func deviceName(name string) string {
	name = strings.TrimSpace(name)
	if utf8.RuneCountInString(name) <= maxDeviceName {
		return name
	}
	return string([]rune(name)[:maxDeviceName])
}
