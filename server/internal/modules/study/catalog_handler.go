package study

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type CatalogHandler struct {
	service Service
}

func (h CatalogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var unitID *string
	if unit := q.Get("unit"); unit != "" {
		unitID = &unit
	}
	list, err := h.service.Catalog(r.Context(), principal(r), unitID, q.Get("archived") == "true")
	if err != nil {
		writeError(w, err)
		return
	}
	if list == nil {
		list = []Course{}
	}
	jsonhttp.Write(w, http.StatusOK, list)
}
