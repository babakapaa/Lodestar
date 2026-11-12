package notes

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type SourcesHandler struct{ service Service }

func (h SourcesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.Sources(r.Context(), principal(r))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, list)
}

type SaveSourceHandler struct{ service Service }

func (h SaveSourceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in SourceInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	src, created, err := h.service.SaveSource(r.Context(), principal(r), in)
	if err != nil {
		writeError(w, err)
		return
	}
	status := http.StatusOK
	if created {
		status = http.StatusCreated
	}
	jsonhttp.Write(w, status, src)
}

type DeleteSourceHandler struct{ service Service }

func (h DeleteSourceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.service.DeleteSource(r.Context(), principal(r), r.PathValue("id")); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type LookupSourceHandler struct{ service Service }

func (h LookupSourceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in LookupInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	fields, err := h.service.LookupSource(r.Context(), in)
	if err != nil {
		writeError(w, err)
		return
	}
	if fields.Authors == nil {
		fields.Authors = []string{}
	}
	jsonhttp.Write(w, http.StatusOK, fields)
}
