package lti

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
)

const (
	assertionType  = "urn:ietf:params:oauth:client-assertion-type:jwt-bearer"
	membersType    = "application/vnd.ims.lti-nrps.v2.membershipcontainer+json"
	lineItemsType  = "application/vnd.ims.lis.v2.lineitemcontainer+json"
	lineItemType   = "application/vnd.ims.lis.v2.lineitem+json"
	scoreType      = "application/vnd.ims.lis.v1.score+json"
	maxPages       = 100
	maxBody        = 8 << 20
	tokenLeeway    = time.Minute
	assertionLife  = 5 * time.Minute
	defaultExpires = time.Hour
)

var nextLink = regexp.MustCompile(`<([^>]+)>\s*;\s*rel="?next"?`)

type PlatformClient struct {
	http   *http.Client
	mu     sync.Mutex
	tokens map[string]serviceToken
}

func NewPlatformClient(client *http.Client) *PlatformClient {
	return &PlatformClient{http: client, tokens: map[string]serviceToken{}}
}

func (c *PlatformClient) Token(ctx context.Context, p Platform, key Key, scopes []string) (string, error) {
	scopes = slices.Clone(scopes)
	slices.Sort(scopes)
	cacheKey := p.ID + " " + strings.Join(scopes, " ")
	c.mu.Lock()
	cached, ok := c.tokens[cacheKey]
	c.mu.Unlock()
	if ok && time.Now().Add(tokenLeeway).Before(cached.expires) {
		return cached.value, nil
	}
	assertion, err := c.assertion(p, key)
	if err != nil {
		return "", err
	}
	form := url.Values{
		"grant_type":            {"client_credentials"},
		"client_assertion_type": {assertionType},
		"client_assertion":      {assertion},
		"scope":                 {strings.Join(scopes, " ")},
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.TokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("lti token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	var out struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := c.do(req, &out); err != nil {
		return "", fmt.Errorf("%w: token: %w", ErrPlatformUnavailable, err)
	}
	if out.AccessToken == "" {
		return "", fmt.Errorf("%w: empty access token", ErrPlatformUnavailable)
	}
	life := defaultExpires
	if out.ExpiresIn > 0 {
		life = time.Duration(out.ExpiresIn) * time.Second
	}
	c.mu.Lock()
	c.tokens[cacheKey] = serviceToken{value: out.AccessToken, expires: time.Now().Add(life)}
	c.mu.Unlock()
	return out.AccessToken, nil
}

func (c *PlatformClient) Members(ctx context.Context, token, membersURL string) ([]NamesMember, error) {
	var all []NamesMember
	err := c.pages(ctx, token, membersURL, membersType, func(body []byte) error {
		var page struct {
			Members []NamesMember `json:"members"`
		}
		if err := json.Unmarshal(body, &page); err != nil {
			return err
		}
		all = append(all, page.Members...)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("%w: members: %w", ErrPlatformUnavailable, err)
	}
	return all, nil
}

func (c *PlatformClient) LineItems(ctx context.Context, token, lineItemsURL string) ([]LineItem, error) {
	var all []LineItem
	err := c.pages(ctx, token, lineItemsURL, lineItemsType, func(body []byte) error {
		var page []LineItem
		if err := json.Unmarshal(body, &page); err != nil {
			return err
		}
		all = append(all, page...)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("%w: line items: %w", ErrPlatformUnavailable, err)
	}
	return all, nil
}

func (c *PlatformClient) CreateLineItem(ctx context.Context, token, lineItemsURL string, item LineItem) (LineItem, error) {
	var out LineItem
	if err := c.send(ctx, token, lineItemsURL, lineItemType, item, &out); err != nil {
		return LineItem{}, fmt.Errorf("%w: create line item: %w", ErrPlatformUnavailable, err)
	}
	if out.ID == "" {
		return LineItem{}, fmt.Errorf("%w: line item without id", ErrPlatformUnavailable)
	}
	return out, nil
}

func (c *PlatformClient) PostScore(ctx context.Context, token, lineItemURL string, s Score) error {
	target, err := url.Parse(lineItemURL)
	if err != nil {
		return fmt.Errorf("lti score url: %w", err)
	}
	target.Path = strings.TrimRight(target.Path, "/") + "/scores"
	if err := c.send(ctx, token, target.String(), scoreType, s, nil); err != nil {
		return fmt.Errorf("%w: score: %w", ErrPlatformUnavailable, err)
	}
	return nil
}

func (c *PlatformClient) assertion(p Platform, key Key) (string, error) {
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: key.Private},
		(&jose.SignerOptions{}).WithType("JWT").WithHeader(jose.HeaderKey("kid"), key.ID))
	if err != nil {
		return "", fmt.Errorf("lti signer: %w", err)
	}
	jti, err := randomString()
	if err != nil {
		return "", err
	}
	now := time.Now()
	claims := jwt.Claims{
		Issuer:   p.ClientID,
		Subject:  p.ClientID,
		Audience: jwt.Audience{p.TokenURL},
		IssuedAt: jwt.NewNumericDate(now),
		Expiry:   jwt.NewNumericDate(now.Add(assertionLife)),
		ID:       jti,
	}
	raw, err := jwt.Signed(signer).Claims(claims).Serialize()
	if err != nil {
		return "", fmt.Errorf("lti assertion: %w", err)
	}
	return raw, nil
}

func (c *PlatformClient) pages(ctx context.Context, token, first, accept string, read func([]byte) error) error {
	next := first
	for page := 0; next != "" && page < maxPages; page++ {
		body, link, err := c.page(ctx, token, next, accept)
		if err != nil {
			return err
		}
		if err := read(body); err != nil {
			return err
		}
		next = link
	}
	return nil
}

func (c *PlatformClient) page(ctx context.Context, token, target, accept string) ([]byte, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", accept)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := readBody(resp)
	if err != nil {
		return nil, "", err
	}
	if m := nextLink.FindStringSubmatch(resp.Header.Get("Link")); m != nil {
		return body, m[1], nil
	}
	return body, "", nil
}

func (c *PlatformClient) send(ctx context.Context, token, target, contentType string, in, out any) error {
	raw, err := json.Marshal(in)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, target, bytes.NewReader(raw))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", contentType)
	if out != nil {
		req.Header.Set("Accept", contentType)
	}
	return c.do(req, out)
}

func (c *PlatformClient) do(req *http.Request, out any) error {
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := readBody(resp)
	if err != nil {
		return err
	}
	if out != nil && len(body) > 0 {
		return json.Unmarshal(body, out)
	}
	return nil
}

func readBody(resp *http.Response) ([]byte, error) {
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBody))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("%s answered %d", resp.Request.URL.Host, resp.StatusCode)
	}
	return body, nil
}
