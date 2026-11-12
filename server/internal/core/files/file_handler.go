package files

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type FileHandler struct {
	service Service
}

func (h FileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f, err := h.service.Get(r.Context(), principal(r), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, f)
}
