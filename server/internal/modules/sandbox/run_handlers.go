package sandbox

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type RunsHandler struct{ service Service }

func (h RunsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.Runs(r.Context(), principal(r))
	respond(w, http.StatusOK, list, err)
}

type StartRunHandler struct{ service Service }

func (h StartRunHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in RunInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	run, err := h.service.StartRun(r.Context(), principal(r), in)
	respond(w, http.StatusCreated, run, err)
}

type RunHandler struct{ service Service }

func (h RunHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	run, err := h.service.GetRun(r.Context(), principal(r), r.PathValue("id"))
	respond(w, http.StatusOK, run, err)
}

type StopRunHandler struct{ service Service }

func (h StopRunHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	run, err := h.service.StopRun(r.Context(), principal(r), r.PathValue("id"))
	respond(w, http.StatusOK, run, err)
}
