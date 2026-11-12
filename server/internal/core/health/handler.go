package health

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type Handler struct{}

func (Handler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	jsonhttp.Write(w, http.StatusOK, Status{Status: StatusOK})
}
