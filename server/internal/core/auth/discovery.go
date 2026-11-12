package auth

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/coreos/go-oidc/v3/oidc"
)

type Discovery struct {
	client    *http.Client
	mu        sync.Mutex
	providers map[string]*oidc.Provider
}

func NewDiscovery(client *http.Client) *Discovery {
	return &Discovery{client: client, providers: map[string]*oidc.Provider{}}
}

func (d *Discovery) Provider(ctx context.Context, issuer string) (*oidc.Provider, error) {
	d.mu.Lock()
	p, ok := d.providers[issuer]
	d.mu.Unlock()
	if ok {
		return p, nil
	}
	p, err := oidc.NewProvider(d.Context(ctx), issuer)
	if err != nil {
		return nil, fmt.Errorf("discover %s: %w", issuer, err)
	}
	d.mu.Lock()
	d.providers[issuer] = p
	d.mu.Unlock()
	return p, nil
}

func (d *Discovery) Context(ctx context.Context) context.Context {
	return oidc.ClientContext(ctx, d.client)
}
