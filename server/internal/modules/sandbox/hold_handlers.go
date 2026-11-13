package sandbox

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type HeldRunsHandler struct{ service Service }

func (h HeldRunsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.HeldRuns(r.Context(), principal(r))
	respond(w, http.StatusOK, list, err)
}

type DecideHoldHandler struct{ service Service }

func (h DecideHoldHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var d HoldDecision
	if err := jsonhttp.Read(w, r, &d); err != nil {
		badRequest(w)
		return
	}
	run, err := h.service.DecideHold(r.Context(), principal(r), r.PathValue("id"), d)
	respond(w, http.StatusOK, run, err)
}
