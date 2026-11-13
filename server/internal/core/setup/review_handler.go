package setup

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type ReviewHandler struct {
	service Service
}

func (h ReviewHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in ReviewInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	sec, err := h.service.Review(r.Context(), principal(r), r.PathValue("section"), in)
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, sec)
}
