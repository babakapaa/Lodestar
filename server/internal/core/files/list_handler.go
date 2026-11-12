package files

import (
	"net/http"
	"time"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type ListHandler struct {
	service Service
}

func (h ListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var before *time.Time
	if value := r.URL.Query().Get("before"); value != "" {
		parsed, err := time.Parse(time.RFC3339Nano, value)
		if err != nil {
			badRequest(w)
			return
		}
		before = &parsed
	}
	list, err := h.service.List(r.Context(), principal(r), before)
	if err != nil {
		writeError(w, err)
		return
	}
	if list == nil {
		list = []File{}
	}
	jsonhttp.Write(w, http.StatusOK, list)
}
