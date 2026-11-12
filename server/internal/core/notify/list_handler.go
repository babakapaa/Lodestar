package notify

import (
	"net/http"
	"time"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type ListHandler struct {
	service Service
}

func (h ListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var before *time.Time
	if value := q.Get("before"); value != "" {
		parsed, err := time.Parse(time.RFC3339Nano, value)
		if err != nil {
			badRequest(w)
			return
		}
		before = &parsed
	}
	list, err := h.service.List(r.Context(), principal(r), before, q.Get("unread") == "true")
	if err != nil {
		writeError(w, err)
		return
	}
	if list == nil {
		list = []Notification{}
	}
	jsonhttp.Write(w, http.StatusOK, list)
}
