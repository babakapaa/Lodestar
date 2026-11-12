package structure

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type UnitsHandler struct {
	service Service
}

func (h UnitsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	list, err := h.service.Children(r.Context(), principal(r), optional(q.Get("parent")), q.Get("archived") == "true")
	if err != nil {
		writeError(w, err)
		return
	}
	if list == nil {
		list = []Unit{}
	}
	jsonhttp.Write(w, http.StatusOK, list)
}
