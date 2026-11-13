package lti

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type LinkHandler struct {
	service Service
}

func (h LinkHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in LinkInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	v, err := h.service.Link(r.Context(), principal(r), r.PathValue("id"), in)
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, v)
}
