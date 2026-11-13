package setup

import (
	"context"
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type StepHandler struct {
	step func(context.Context, auth.Principal) (Overview, error)
}

func (h StepHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	o, err := h.step(r.Context(), principal(r))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, o)
}
