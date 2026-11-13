package sandbox

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type EstimateHandler struct{ service Service }

func (h EstimateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in EstimateInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	e, err := h.service.Estimate(r.Context(), principal(r), in)
	respond(w, http.StatusOK, e, err)
}
