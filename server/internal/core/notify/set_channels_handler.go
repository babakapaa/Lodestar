package notify

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type SetChannelsHandler struct {
	service Service
}

func (h SetChannelsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in ChannelsInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	if err := h.service.SetChannels(r.Context(), principal(r), r.PathValue("type"), in); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
