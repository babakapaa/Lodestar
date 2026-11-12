package files

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type DownloadHandler struct {
	service Service
}

func (h DownloadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d, err := h.service.Download(r.Context(), principal(r), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, d)
}
