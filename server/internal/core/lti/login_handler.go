package lti

import (
	"log/slog"
	"net/http"
)

type LoginHandler struct {
	launch LaunchService
}

func (h LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		badRequest(w)
		return
	}
	target, err := h.launch.Begin(r.Context(), LoginInit{
		TenantSlug:   r.PathValue("tenant"),
		Issuer:       r.Form.Get("iss"),
		ClientID:     r.Form.Get("client_id"),
		DeploymentID: r.Form.Get("lti_deployment_id"),
		LoginHint:    r.Form.Get("login_hint"),
		MessageHint:  r.Form.Get("lti_message_hint"),
	})
	if err != nil {
		slog.Warn("lti login rejected", "tenant", r.PathValue("tenant"), "issuer", r.Form.Get("iss"), "error", err)
		writeError(w, err)
		return
	}
	http.Redirect(w, r, target, http.StatusFound)
}
