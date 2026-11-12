package settings

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type SchemaHandler struct {
	service Service
}

func (h SchemaHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	jsonhttp.Write(w, http.StatusOK, h.service.Catalog())
}
