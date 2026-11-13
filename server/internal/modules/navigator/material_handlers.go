package navigator

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

const materialBody = 4 << 20

type MaterialsHandler struct{ service MaterialService }

func (h MaterialsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.List(r.Context(), principal(r), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, list)
}

type AddMaterialHandler struct{ service MaterialService }

func (h AddMaterialHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in MaterialInput
	if err := jsonhttp.ReadLimit(w, r, &in, materialBody); err != nil {
		badRequest(w)
		return
	}
	m, err := h.service.Add(r.Context(), principal(r), r.PathValue("id"), in)
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusCreated, m)
}

type DeleteMaterialHandler struct{ service MaterialService }

func (h DeleteMaterialHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Delete(r.Context(), principal(r), r.PathValue("id")); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
