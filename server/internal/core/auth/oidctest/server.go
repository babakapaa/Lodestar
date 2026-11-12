package oidctest

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
)

const keyID = "test-key"

type Server struct {
	URL          string
	ClientID     string
	ClientSecret string

	t      *testing.T
	key    *rsa.PrivateKey
	signer jose.Signer
	mu     sync.Mutex
	user   User
	codes  map[string]pendingCode
}

func New(t *testing.T) *Server {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	signer, err := jose.NewSigner(
		jose.SigningKey{Algorithm: jose.RS256, Key: key},
		(&jose.SignerOptions{}).WithType("JWT").WithHeader(jose.HeaderKey("kid"), keyID))
	if err != nil {
		t.Fatalf("signer: %v", err)
	}
	s := &Server{
		ClientID:     "lodestar",
		ClientSecret: "test-secret",
		t:            t,
		key:          key,
		signer:       signer,
		user:         User{Subject: "student-1", Email: "student@uni.example", Name: "Айгерим"},
		codes:        map[string]pendingCode{},
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /.well-known/openid-configuration", s.discovery)
	mux.HandleFunc("GET /authorize", s.authorize)
	mux.HandleFunc("POST /token", s.token)
	mux.HandleFunc("GET /jwks", s.jwks)
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	s.URL = srv.URL
	return s
}

func (s *Server) SignIn(u User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.user = u
}

func (s *Server) discovery(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, map[string]any{
		"issuer":                                s.URL,
		"authorization_endpoint":                s.URL + "/authorize",
		"token_endpoint":                        s.URL + "/token",
		"jwks_uri":                              s.URL + "/jwks",
		"response_types_supported":              []string{"code"},
		"subject_types_supported":               []string{"public"},
		"id_token_signing_alg_values_supported": []string{"RS256"},
	})
}

func (s *Server) authorize(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	if q.Get("client_id") != s.ClientID || q.Get("code_challenge_method") != "S256" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	code := rand.Text()
	s.mu.Lock()
	s.codes[code] = pendingCode{
		clientID:    q.Get("client_id"),
		redirectURI: q.Get("redirect_uri"),
		nonce:       q.Get("nonce"),
		challenge:   q.Get("code_challenge"),
		user:        s.user,
	}
	s.mu.Unlock()
	target, err := url.Parse(q.Get("redirect_uri"))
	if err != nil {
		http.Error(w, "bad redirect", http.StatusBadRequest)
		return
	}
	back := target.Query()
	back.Set("code", code)
	back.Set("state", q.Get("state"))
	target.RawQuery = back.Encode()
	http.Redirect(w, r, target.String(), http.StatusFound)
}

func (s *Server) token(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	id, secret, ok := r.BasicAuth()
	if !ok {
		id, secret = r.PostForm.Get("client_id"), r.PostForm.Get("client_secret")
	}
	s.mu.Lock()
	pending, found := s.codes[r.PostForm.Get("code")]
	delete(s.codes, r.PostForm.Get("code"))
	s.mu.Unlock()
	if !found || id != s.ClientID || secret != s.ClientSecret ||
		pending.redirectURI != r.PostForm.Get("redirect_uri") ||
		!challengeMatches(pending.challenge, r.PostForm.Get("code_verifier")) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid_grant"})
		return
	}
	idToken, err := s.idToken(pending)
	if err != nil {
		s.t.Errorf("sign id_token: %v", err)
		http.Error(w, "sign", http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]any{
		"access_token": rand.Text(),
		"token_type":   "Bearer",
		"expires_in":   300,
		"id_token":     idToken,
	})
}

func (s *Server) idToken(p pendingCode) (string, error) {
	now := time.Now()
	claims := struct {
		jwt.Claims
		Nonce string `json:"nonce"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}{
		Claims: jwt.Claims{
			Issuer:   s.URL,
			Subject:  p.user.Subject,
			Audience: jwt.Audience{p.clientID},
			IssuedAt: jwt.NewNumericDate(now),
			Expiry:   jwt.NewNumericDate(now.Add(5 * time.Minute)),
		},
		Nonce: p.nonce,
		Email: p.user.Email,
		Name:  p.user.Name,
	}
	return jwt.Signed(s.signer).Claims(claims).Serialize()
}

func (s *Server) jwks(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, jose.JSONWebKeySet{Keys: []jose.JSONWebKey{{
		Key:       &s.key.PublicKey,
		KeyID:     keyID,
		Algorithm: string(jose.RS256),
		Use:       "sig",
	}}})
}

func challengeMatches(challenge, verifier string) bool {
	sum := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(sum[:]) == challenge
}

func writeJSON(w http.ResponseWriter, body any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(body)
}
