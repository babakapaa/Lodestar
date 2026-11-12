package privacy

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type DeletionHandler struct {
	service Service
}

func (h DeletionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d, err := h.service.Deletion(r.Context(), principal(r))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, d)
}
