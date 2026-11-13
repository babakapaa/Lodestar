package navigator

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type DialogsHandler struct{ service Service }

func (h DialogsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	list, err := h.service.Dialogs(r.Context(), principal(r), q.Get("scopeKind"), q.Get("scopeId"))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, list)
}

type CreateHandler struct{ service Service }

func (h CreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in DialogInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	d, created, err := h.service.Create(r.Context(), principal(r), in)
	if err != nil {
		writeError(w, err)
		return
	}
	status := http.StatusOK
	if created {
		status = http.StatusCreated
	}
	jsonhttp.Write(w, status, d)
}

type DialogHandler struct{ service Service }

func (h DialogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d, err := h.service.Get(r.Context(), principal(r), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, d)
}

type DeleteHandler struct{ service Service }

func (h DeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Delete(r.Context(), principal(r), r.PathValue("id")); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
