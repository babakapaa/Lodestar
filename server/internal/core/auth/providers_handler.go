package auth

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type ProvidersHandler struct {
	login LoginService
}

func (h ProvidersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list, err := h.login.Providers(r.Context(), r.PathValue("tenant"))
	if err != nil {
		writeError(w, err)
		return
	}
	if list == nil {
		list = []ProviderInfo{}
	}
	jsonhttp.Write(w, http.StatusOK, list)
}
