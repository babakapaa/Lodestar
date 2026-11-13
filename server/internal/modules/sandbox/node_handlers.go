package sandbox

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type NodesHandler struct{ service Service }

func (h NodesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.Nodes(r.Context(), principal(r))
	respond(w, http.StatusOK, list, err)
}

type CreateNodeHandler struct{ service Service }

func (h CreateNodeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in NodeInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	n, err := h.service.CreateNode(r.Context(), principal(r), in)
	respond(w, http.StatusCreated, n, err)
}

type DeleteNodeHandler struct{ service Service }

func (h DeleteNodeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.service.DeleteNode(r.Context(), principal(r), r.PathValue("id")); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
