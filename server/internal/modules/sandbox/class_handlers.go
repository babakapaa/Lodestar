package sandbox

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type ClassesHandler struct{ service Service }

func (h ClassesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.Classes(r.Context(), principal(r))
	respond(w, http.StatusOK, list, err)
}

type CreateClassHandler struct{ service Service }

func (h CreateClassHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in GPUClassInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	c, err := h.service.CreateClass(r.Context(), principal(r), in)
	respond(w, http.StatusCreated, c, err)
}

type UpdateClassHandler struct{ service Service }

func (h UpdateClassHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in GPUClassInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	c, err := h.service.UpdateClass(r.Context(), principal(r), r.PathValue("id"), in)
	respond(w, http.StatusOK, c, err)
}
