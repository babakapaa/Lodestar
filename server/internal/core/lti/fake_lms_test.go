package lti_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/babakapaa/lodestar/server/internal/core/lti"
)

type fakeLMS struct {
	mu        sync.Mutex
	url       string
	members   []lti.NamesMember
	items     []lti.LineItem
	scores    []lti.Score
	tokens    int
	assertion string
}

func (f *fakeLMS) handle(w http.ResponseWriter, r *http.Request) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	switch {
	case r.URL.Path == "/token":
		_ = r.ParseForm()
		if r.PostForm.Get("grant_type") != "client_credentials" || r.PostForm.Get("client_assertion") == "" {
			http.Error(w, "bad", http.StatusBadRequest)
			return true
		}
		f.tokens++
		f.assertion = r.PostForm.Get("client_assertion")
		_ = json.NewEncoder(w).Encode(map[string]any{"access_token": "svc", "expires_in": 3600})
	case r.Header.Get("Authorization") != "Bearer svc":
		http.Error(w, "no", http.StatusUnauthorized)
	case r.URL.Path == "/memberships":
		_ = json.NewEncoder(w).Encode(map[string]any{"members": f.members})
	case r.URL.Path == "/lineitems" && r.Method == http.MethodGet:
		_ = json.NewEncoder(w).Encode(f.items)
	case r.URL.Path == "/lineitems" && r.Method == http.MethodPost:
		var item lti.LineItem
		_ = json.NewDecoder(r.Body).Decode(&item)
		item.ID = f.url + "/lineitems/" + item.ResourceID
		f.items = append(f.items, item)
		_ = json.NewEncoder(w).Encode(item)
	case strings.HasSuffix(r.URL.Path, "/scores"):
		var s lti.Score
		_ = json.NewDecoder(r.Body).Decode(&s)
		f.scores = append(f.scores, s)
		w.WriteHeader(http.StatusNoContent)
	default:
		return false
	}
	return true
}
