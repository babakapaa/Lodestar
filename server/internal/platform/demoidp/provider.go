package demoidp

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"slices"
	"sync"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
)

const (
	keyID    = "demo-key"
	codeTTL  = 5 * time.Minute
	tokenTTL = 5 * time.Minute
)

var ErrNoUsers = errors.New("demo identity provider needs users or a pick function")

type Provider struct {
	cfg    Config
	key    *rsa.PrivateKey
	signer jose.Signer
	mu     sync.Mutex
	codes  map[string]pendingCode
}

func New(cfg Config) (*Provider, error) {
	if len(cfg.Users) == 0 && cfg.Pick == nil {
		return nil, ErrNoUsers
	}
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	signer, err := jose.NewSigner(
		jose.SigningKey{Algorithm: jose.RS256, Key: key},
		(&jose.SignerOptions{}).WithType("JWT").WithHeader(jose.HeaderKey("kid"), keyID))
	if err != nil {
		return nil, err
	}
	return &Provider{cfg: cfg, key: key, signer: signer, codes: map[string]pendingCode{}}, nil
}

func (p *Provider) SetIssuer(issuer string) {
	p.cfg.Issuer = issuer
}

func (p *Provider) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /.well-known/openid-configuration", p.discovery)
	mux.HandleFunc("GET /authorize", p.authorize)
	mux.HandleFunc("POST /authorize", p.choose)
	mux.HandleFunc("POST /token", p.token)
	mux.HandleFunc("GET /jwks", p.jwks)
	return mux
}

func (p *Provider) discovery(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"issuer":                                p.cfg.Issuer,
		"authorization_endpoint":                p.cfg.Issuer + "/authorize",
		"token_endpoint":                        p.cfg.Issuer + "/token",
		"jwks_uri":                              p.cfg.Issuer + "/jwks",
		"response_types_supported":              []string{"code"},
		"subject_types_supported":               []string{"public"},
		"id_token_signing_alg_values_supported": []string{"RS256"},
	})
}

func (p *Provider) authorize(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	if !p.validRequest(q) {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if p.cfg.Pick != nil {
		p.redirect(w, r, q, p.cfg.Pick())
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = choicePage.Execute(w, choiceView{Users: p.cfg.Users, Query: q})
}

func (p *Provider) choose(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil || !p.validRequest(r.PostForm) {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	subject := r.PostForm.Get("subject")
	i := slices.IndexFunc(p.cfg.Users, func(u User) bool { return u.Subject == subject })
	if i < 0 {
		http.Error(w, "unknown user", http.StatusBadRequest)
		return
	}
	p.redirect(w, r, r.PostForm, p.cfg.Users[i])
}

func (p *Provider) validRequest(q url.Values) bool {
	if q.Get("client_id") != p.cfg.ClientID || q.Get("code_challenge_method") != "S256" || q.Get("code_challenge") == "" {
		return false
	}
	target, err := url.Parse(q.Get("redirect_uri"))
	if err != nil || target.Scheme == "" {
		return false
	}
	return len(p.cfg.Redirects) == 0 || slices.Contains(p.cfg.Redirects, q.Get("redirect_uri"))
}

func (p *Provider) redirect(w http.ResponseWriter, r *http.Request, q url.Values, u User) {
	code := rand.Text()
	p.mu.Lock()
	p.codes[code] = pendingCode{
		redirectURI: q.Get("redirect_uri"),
		nonce:       q.Get("nonce"),
		challenge:   q.Get("code_challenge"),
		user:        u,
	}
	p.mu.Unlock()
	time.AfterFunc(codeTTL, func() {
		p.mu.Lock()
		delete(p.codes, code)
		p.mu.Unlock()
	})
	target, _ := url.Parse(q.Get("redirect_uri"))
	back := target.Query()
	back.Set("code", code)
	back.Set("state", q.Get("state"))
	target.RawQuery = back.Encode()
	http.Redirect(w, r, target.String(), http.StatusFound)
}

func (p *Provider) token(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	id, secret, ok := r.BasicAuth()
	if !ok {
		id, secret = r.PostForm.Get("client_id"), r.PostForm.Get("client_secret")
	}
	code := r.PostForm.Get("code")
	p.mu.Lock()
	pending, found := p.codes[code]
	delete(p.codes, code)
	p.mu.Unlock()
	if !found || id != p.cfg.ClientID || secret != p.cfg.ClientSecret ||
		pending.redirectURI != r.PostForm.Get("redirect_uri") ||
		!challengeMatches(pending.challenge, r.PostForm.Get("code_verifier")) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid_grant"})
		return
	}
	idToken, err := p.idToken(id, pending)
	if err != nil {
		http.Error(w, "sign", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"access_token": rand.Text(),
		"token_type":   "Bearer",
		"expires_in":   int(tokenTTL.Seconds()),
		"id_token":     idToken,
	})
}

func (p *Provider) idToken(clientID string, pending pendingCode) (string, error) {
	now := time.Now()
	claims := struct {
		jwt.Claims
		Nonce string `json:"nonce"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}{
		Claims: jwt.Claims{
			Issuer:   p.cfg.Issuer,
			Subject:  pending.user.Subject,
			Audience: jwt.Audience{clientID},
			IssuedAt: jwt.NewNumericDate(now),
			Expiry:   jwt.NewNumericDate(now.Add(tokenTTL)),
		},
		Nonce: pending.nonce,
		Email: pending.user.Email,
		Name:  pending.user.Name,
	}
	return jwt.Signed(p.signer).Claims(claims).Serialize()
}

func (p *Provider) jwks(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, jose.JSONWebKeySet{Keys: []jose.JSONWebKey{{
		Key:       &p.key.PublicKey,
		KeyID:     keyID,
		Algorithm: string(jose.RS256),
		Use:       "sig",
	}}})
}

func challengeMatches(challenge, verifier string) bool {
	sum := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(sum[:]) == challenge
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
