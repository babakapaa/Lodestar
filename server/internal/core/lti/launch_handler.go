package lti

import (
	"errors"
	"log/slog"
	"net/http"
)

type LaunchHandler struct {
	launch LaunchService
}

func (h LaunchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		badRequest(w)
		return
	}
	target, err := h.launch.Complete(r.Context(), Launch{
		TenantSlug: r.PathValue("tenant"),
		State:      r.PostForm.Get("state"),
		IDToken:    r.PostForm.Get("id_token"),
		Error:      r.PostForm.Get("error"),
	})
	if errors.Is(err, ErrLaunchRejected) {
		slog.Warn("lti launch rejected", "tenant", r.PathValue("tenant"), "error", err)
		http.Redirect(w, r, target, http.StatusSeeOther)
		return
	}
	if err != nil {
		writeError(w, err)
		return
	}
	http.Redirect(w, r, target, http.StatusSeeOther)
}
