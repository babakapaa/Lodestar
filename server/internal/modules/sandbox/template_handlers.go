package sandbox

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type TemplatesHandler struct{ service Service }

func (h TemplatesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.Templates(r.Context(), principal(r))
	respond(w, http.StatusOK, list, err)
}

type CreateTemplateHandler struct{ service Service }

func (h CreateTemplateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in TemplateInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	t, err := h.service.CreateTemplate(r.Context(), principal(r), in)
	respond(w, http.StatusCreated, t, err)
}

type UpdateTemplateHandler struct{ service Service }

func (h UpdateTemplateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in TemplateInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	t, err := h.service.UpdateTemplate(r.Context(), principal(r), r.PathValue("id"), in)
	respond(w, http.StatusOK, t, err)
}
