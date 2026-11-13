package lti

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/core/tenant"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

const appLaunchPath = "/lti"

type LaunchService struct {
	tenants   tenant.Store
	tx        database.TenantTx
	accounts  Accounts
	tickets   auth.TicketStore
	keySets   *KeySets
	links     module.LinkSources
	platforms PlatformStore
	attempts  AttemptStore
	contexts  ContextStore
	members   MemberStore
	publicURL string
	appURL    string
}

func (s LaunchService) Begin(ctx context.Context, in LoginInit) (string, error) {
	t, err := s.tenants.BySlug(ctx, in.TenantSlug)
	if err != nil {
		return "", err
	}
	var p Platform
	attempt := Attempt{}
	if attempt.State, err = randomString(); err != nil {
		return "", err
	}
	if attempt.Nonce, err = randomString(); err != nil {
		return "", err
	}
	err = s.tx.Run(ctx, t.ID, func(tx pgx.Tx) error {
		var err error
		if p, err = s.platforms.ByIssuer(ctx, tx, in.Issuer, in.ClientID); err != nil {
			return err
		}
		if in.DeploymentID != "" && !p.Deploys(in.DeploymentID) {
			return ErrPlatformNotFound
		}
		attempt.PlatformID = p.ID
		return s.attempts.Create(ctx, tx, attempt)
	})
	if err != nil {
		return "", err
	}
	target, err := url.Parse(p.AuthURL)
	if err != nil {
		return "", ErrInvalidPlatform
	}
	q := target.Query()
	q.Set("scope", "openid")
	q.Set("response_type", "id_token")
	q.Set("response_mode", "form_post")
	q.Set("prompt", "none")
	q.Set("client_id", p.ClientID)
	q.Set("redirect_uri", NewToolURLs(s.publicURL, t.Slug).LaunchURL)
	q.Set("login_hint", in.LoginHint)
	q.Set("state", attempt.State)
	q.Set("nonce", attempt.Nonce)
	if in.MessageHint != "" {
		q.Set("lti_message_hint", in.MessageHint)
	}
	target.RawQuery = q.Encode()
	return target.String(), nil
}

func (s LaunchService) Complete(ctx context.Context, in Launch) (string, error) {
	t, err := s.tenants.BySlug(ctx, in.TenantSlug)
	if err != nil {
		return "", err
	}
	next, err := s.complete(ctx, t, in)
	if err != nil {
		return s.appLink("error", "launch_failed"), fmt.Errorf("%w: %w", ErrLaunchRejected, err)
	}
	return next, nil
}

func (s LaunchService) complete(ctx context.Context, t tenant.Tenant, in Launch) (string, error) {
	if in.Error != "" {
		return "", fmt.Errorf("platform returned %s", in.Error)
	}
	var p Platform
	var attempt Attempt
	err := s.tx.Run(ctx, t.ID, func(tx pgx.Tx) error {
		var err error
		if attempt, err = s.attempts.Take(ctx, tx, in.State); err != nil {
			return err
		}
		p, err = s.platforms.ByID(ctx, tx, attempt.PlatformID)
		return err
	})
	if err != nil {
		return "", err
	}
	token, err := s.keySets.Verifier(p).Verify(ctx, in.IDToken)
	if err != nil {
		return "", fmt.Errorf("verify id_token: %w", err)
	}
	if token.Nonce != attempt.Nonce {
		return "", fmt.Errorf("id_token nonce mismatch")
	}
	var claims LaunchClaims
	if err := token.Claims(&claims); err != nil {
		return "", fmt.Errorf("read claims: %w", err)
	}
	if claims.AuthorizedBy != "" && claims.AuthorizedBy != p.ClientID {
		return "", fmt.Errorf("id_token azp mismatch")
	}
	if !claims.Valid() || !p.Deploys(claims.DeploymentID) {
		return "", fmt.Errorf("unsupported launch %q for deployment %q", claims.MessageType, claims.DeploymentID)
	}
	userID, err := s.accounts.Resolve(ctx, t.ID, p, Person{Subject: token.Subject, Email: claims.Email, Name: claims.DisplayName()})
	if err != nil {
		return "", err
	}
	var path string
	err = s.tx.Run(ctx, t.ID, func(tx pgx.Tx) error {
		c, err := s.contexts.Upsert(ctx, tx, p.ID, claims)
		if err != nil {
			return err
		}
		if err := s.members.Record(ctx, tx, c.ID, userID, Member{Roles: claims.Roles, Instructor: claims.Instructor()}); err != nil {
			return err
		}
		path, err = s.path(ctx, tx, c)
		return err
	})
	if err != nil {
		return "", err
	}
	ticket, err := s.tickets.Create(ctx, t.ID, userID, p.Name)
	if err != nil {
		return "", err
	}
	return s.appLink("ticket", ticket, "next", path, "server", s.publicURL, "tenant", t.Slug), nil
}

func (s LaunchService) path(ctx context.Context, tx pgx.Tx, c Context) (string, error) {
	fallback := appLaunchPath + "/contexts/" + c.ID
	if c.OriginKind == nil {
		return fallback, nil
	}
	source, ok := s.links.Find(*c.OriginKind)
	if !ok {
		return fallback, nil
	}
	target, err := source.Describe(ctx, tx, *c.OriginID)
	if errors.Is(err, module.ErrLinkTargetNotFound) {
		return fallback, nil
	}
	if err != nil {
		return "", err
	}
	return target.Path, nil
}

func (s LaunchService) appLink(pairs ...string) string {
	q := url.Values{}
	for i := 0; i+1 < len(pairs); i += 2 {
		q.Set(pairs[i], pairs[i+1])
	}
	return s.appURL + appLaunchPath + "?" + q.Encode()
}
