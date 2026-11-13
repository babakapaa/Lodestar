package auth

import (
	"encoding/json"
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/apierror"
	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

const maxTokenRequest = 8 << 10

type TokenHandler struct {
	login    LoginService
	sessions SessionStore
}

func (h TokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req TokenRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxTokenRequest)).Decode(&req); err != nil {
		apierror.Write(w, http.StatusBadRequest, CodeInvalidRequest)
		return
	}
	var (
		pair TokenPair
		err  error
	)
	switch req.GrantType {
	case "authorization_code":
		pair, err = h.login.Exchange(r.Context(), req.Code, req.CodeVerifier, req.RedirectURI)
	case "launch_ticket":
		pair, err = h.login.Claim(r.Context(), req.Ticket)
	case "refresh_token":
		pair, err = h.sessions.Refresh(r.Context(), req.RefreshToken)
	default:
		apierror.Write(w, http.StatusBadRequest, CodeInvalidRequest)
		return
	}
	if err != nil {
		writeError(w, err)
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	jsonhttp.Write(w, http.StatusOK, pair)
}
