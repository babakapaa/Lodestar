package files

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type CompleteHandler struct {
	service Service
}

func (h CompleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f, err := h.service.Complete(r.Context(), principal(r), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, f)
}
