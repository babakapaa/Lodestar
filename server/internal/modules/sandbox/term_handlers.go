package sandbox

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type TermsHandler struct{ service Service }

func (h TermsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.Terms(r.Context(), principal(r))
	respond(w, http.StatusOK, list, err)
}

type CreateTermHandler struct{ service Service }

func (h CreateTermHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in TermInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	t, err := h.service.CreateTerm(r.Context(), principal(r), in)
	respond(w, http.StatusCreated, t, err)
}

type UpdateTermHandler struct{ service Service }

func (h UpdateTermHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in TermInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	t, err := h.service.UpdateTerm(r.Context(), principal(r), r.PathValue("id"), in)
	respond(w, http.StatusOK, t, err)
}

type DeleteTermHandler struct{ service Service }

func (h DeleteTermHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.service.DeleteTerm(r.Context(), principal(r), r.PathValue("id")); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
