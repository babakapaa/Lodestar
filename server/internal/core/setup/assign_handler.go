package setup

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type AssignHandler struct {
	service Service
}

func (h AssignHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in AssignInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	sec, err := h.service.Assign(r.Context(), principal(r), r.PathValue("section"), in)
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, sec)
}
