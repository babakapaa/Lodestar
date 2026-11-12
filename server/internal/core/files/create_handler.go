package files

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type CreateHandler struct {
	service Service
}

func (h CreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in FileInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	created, err := h.service.Create(r.Context(), principal(r), in)
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusCreated, created)
}
