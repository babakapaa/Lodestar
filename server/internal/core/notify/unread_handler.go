package notify

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type UnreadHandler struct {
	service Service
}

func (h UnreadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	n, err := h.service.Unread(r.Context(), principal(r))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, UnreadCount{Unread: n})
}
