package study

import (
	"net/http"
	"time"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type PreviewHandler struct {
	service AssignmentService
}

func (h PreviewHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	at := time.Now()
	if raw := r.URL.Query().Get("at"); raw != "" {
		parsed, err := time.Parse(time.RFC3339, raw)
		if err != nil {
			badRequest(w)
			return
		}
		at = parsed
	}
	list, err := h.service.Preview(r.Context(), principal(r), r.PathValue("id"), r.URL.Query().Get("student"), at)
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, list)
}
