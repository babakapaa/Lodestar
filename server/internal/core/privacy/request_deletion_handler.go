package privacy

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type RequestDeletionHandler struct {
	service Service
}

func (h RequestDeletionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d, err := h.service.RequestDeletion(r.Context(), principal(r))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, d)
}
