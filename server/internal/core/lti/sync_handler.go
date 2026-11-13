package lti

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type SyncHandler struct {
	service Service
}

func (h SyncHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	report, err := h.service.Sync(r.Context(), principal(r), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, report)
}
