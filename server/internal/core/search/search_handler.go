package search

import (
	"net/http"
	"strconv"

	"github.com/babakapaa/lodestar/server/internal/core/apierror"
	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type SearchHandler struct {
	service Service
}

func (h SearchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit := 0
	if value := q.Get("limit"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			apierror.Write(w, http.StatusBadRequest, CodeInvalidQuery)
			return
		}
		limit = parsed
	}
	hits, err := h.service.Search(r.Context(), principal(r), q.Get("q"), q["kind"], limit)
	if err != nil {
		writeError(w, err)
		return
	}
	if hits == nil {
		hits = []Hit{}
	}
	jsonhttp.Write(w, http.StatusOK, hits)
}
