package delivery

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/apierror"
	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type RegisterHandler struct {
	service Service
}

func (h RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in DeviceInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		apierror.Write(w, http.StatusBadRequest, CodeInvalidRequest)
		return
	}
	d, err := h.service.Register(r.Context(), principal(r), in)
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, d)
}
