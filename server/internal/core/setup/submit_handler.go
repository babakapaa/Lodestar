package setup

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type SubmitHandler struct {
	service Service
}

func (h SubmitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sec, err := h.service.Submit(r.Context(), principal(r), r.PathValue("section"))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, sec)
}
