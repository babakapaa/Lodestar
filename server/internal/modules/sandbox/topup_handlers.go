package sandbox

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type TopupsHandler struct{ service Service }

func (h TopupsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.Topups(r.Context(), principal(r))
	respond(w, http.StatusOK, list, err)
}

type RequestTopupHandler struct{ service Service }

func (h RequestTopupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in TopupInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	t, err := h.service.RequestTopup(r.Context(), principal(r), in)
	respond(w, http.StatusCreated, t, err)
}

type DecideTopupHandler struct{ service Service }

func (h DecideTopupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var d TopupDecision
	if err := jsonhttp.Read(w, r, &d); err != nil {
		badRequest(w)
		return
	}
	t, err := h.service.DecideTopup(r.Context(), principal(r), r.PathValue("id"), d)
	respond(w, http.StatusOK, t, err)
}
