package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/babakapaa/lodestar/server/internal/core/apierror"
)

type Authenticator struct {
	sessions SessionStore
}

func NewAuthenticator(sessions SessionStore) Authenticator {
	return Authenticator{sessions: sessions}
}

func (a Authenticator) Require(next http.Handler) http.Handler {
	return a.require(next, false)
}

func (a Authenticator) RequireStream(next http.Handler) http.Handler {
	return a.require(next, true)
}

func (a Authenticator) require(next http.Handler, allowQuery bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		value, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
		if !ok && allowQuery {
			value = r.URL.Query().Get("access_token")
			ok = value != ""
		}
		if !ok {
			apierror.Write(w, http.StatusUnauthorized, CodeUnauthorized)
			return
		}
		p, err := a.sessions.Authenticate(r.Context(), value)
		if errors.Is(err, ErrInvalidToken) {
			apierror.Write(w, http.StatusUnauthorized, CodeUnauthorized)
			return
		}
		if err != nil {
			writeError(w, err)
			return
		}
		next.ServeHTTP(w, r.WithContext(WithPrincipal(r.Context(), p)))
	})
}
