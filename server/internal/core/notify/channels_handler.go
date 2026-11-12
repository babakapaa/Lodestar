package notify

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type ChannelsHandler struct {
	service Service
}

func (h ChannelsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.Channels(r.Context(), principal(r))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, list)
}
