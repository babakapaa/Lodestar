package auth

import "net/http"

type LoginHandler struct {
	login LoginService
}

func (h LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	if method := q.Get("code_challenge_method"); method != "S256" {
		writeError(w, ErrInvalidChallenge)
		return
	}
	target, err := h.login.Begin(r.Context(), LoginRequest{
		TenantSlug:    r.PathValue("tenant"),
		ProviderKey:   r.PathValue("provider"),
		RedirectURI:   q.Get("redirect_uri"),
		CodeChallenge: q.Get("code_challenge"),
		DeviceName:    q.Get("device"),
	})
	if err != nil {
		writeError(w, err)
		return
	}
	http.Redirect(w, r, target, http.StatusFound)
}
