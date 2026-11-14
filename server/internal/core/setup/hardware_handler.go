package setup

import (
	"context"
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type HardwareHandler struct {
	check func(context.Context, auth.Principal) (HardwareReport, error)
}

func (h HardwareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	report, err := h.check(r.Context(), principal(r))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, report)
}
