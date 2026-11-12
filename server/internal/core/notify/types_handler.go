package notify

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type TypesHandler struct {
	service Service
}

func (h TypesHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	jsonhttp.Write(w, http.StatusOK, h.service.Types())
}
