package sandbox

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type QuotaHandler struct{ service Service }

func (h QuotaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q, err := h.service.Quota(r.Context(), principal(r))
	respond(w, http.StatusOK, q, err)
}

type TransferHandler struct{ service Service }

func (h TransferHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in TransferInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	q, err := h.service.Transfer(r.Context(), principal(r), in)
	respond(w, http.StatusOK, q, err)
}
