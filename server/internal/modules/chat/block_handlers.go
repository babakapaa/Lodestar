package chat

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type BlocksHandler struct{ service BlockService }

func (h BlocksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.List(r.Context(), principal(r))
	writeOK(w, list, err)
}

type BlockHandler struct{ service BlockService }

func (h BlockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in DirectInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	if err := h.service.Block(r.Context(), principal(r), in); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type UnblockHandler struct{ service BlockService }

func (h UnblockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Unblock(r.Context(), principal(r), r.PathValue("userId")); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
