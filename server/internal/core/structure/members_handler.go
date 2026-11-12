package structure

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type MembersHandler struct {
	service Service
}

func (h MembersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.Members(r.Context(), principal(r), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	if list == nil {
		list = []Member{}
	}
	jsonhttp.Write(w, http.StatusOK, list)
}
