package auth

import (
	"errors"
	"log/slog"
	"net/http"
)

type CallbackHandler struct {
	login LoginService
}

func (h CallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	target, err := h.login.Complete(r.Context(), Callback{
		TenantSlug:  r.PathValue("tenant"),
		ProviderKey: r.PathValue("provider"),
		State:       q.Get("state"),
		Code:        q.Get("code"),
		Error:       q.Get("error"),
	})
	if errors.Is(err, ErrLoginFailed) {
		slog.Warn("login failed", "tenant", r.PathValue("tenant"), "provider", r.PathValue("provider"), "error", err)
		http.Redirect(w, r, target, http.StatusFound)
		return
	}
	if err != nil {
		writeError(w, err)
		return
	}
	http.Redirect(w, r, target, http.StatusFound)
}
