package profile

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type PreferenceSchemaHandler struct {
	service Service
}

func (h PreferenceSchemaHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	jsonhttp.Write(w, http.StatusOK, h.service.Catalog())
}
