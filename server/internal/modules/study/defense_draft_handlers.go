package study

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type DraftProtocolHandler struct{ service DefenseService }

func (h DraftProtocolHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in DefenseInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	summary, err := h.service.DraftProtocol(r.Context(), principal(r), r.PathValue("id"), r.PathValue("student"), in)
	respond(w, http.StatusOK, map[string]any{"summary": summary}, err)
}
