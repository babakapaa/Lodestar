package structure

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type UserUnitsHandler struct {
	service Service
}

func (h UserUnitsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.UserUnits(r.Context(), principal(r), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	if list == nil {
		list = []Unit{}
	}
	jsonhttp.Write(w, http.StatusOK, list)
}
