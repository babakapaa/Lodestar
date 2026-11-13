package lti

import (
	"context"
	"net/http"
	"sync"

	"github.com/coreos/go-oidc/v3/oidc"
)

type KeySets struct {
	client *http.Client
	mu     sync.Mutex
	sets   map[string]*oidc.RemoteKeySet
}

func NewKeySets(client *http.Client) *KeySets {
	return &KeySets{client: client, sets: map[string]*oidc.RemoteKeySet{}}
}

func (k *KeySets) Verifier(p Platform) *oidc.IDTokenVerifier {
	k.mu.Lock()
	set, ok := k.sets[p.JWKSURL]
	if !ok {
		set = oidc.NewRemoteKeySet(oidc.ClientContext(context.Background(), k.client), p.JWKSURL)
		k.sets[p.JWKSURL] = set
	}
	k.mu.Unlock()
	return oidc.NewVerifier(p.Issuer, set, &oidc.Config{ClientID: p.ClientID})
}
