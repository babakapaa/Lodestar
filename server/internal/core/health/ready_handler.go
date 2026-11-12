package health

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/babakapaa/lodestar/server/internal/core/apierror"
	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

const CodeNotReady = "not_ready"

type ReadyHandler struct {
	db Pinger
}

func NewReadyHandler(db Pinger) ReadyHandler {
	return ReadyHandler{db: db}
}

func (h ReadyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	if err := h.db.Ping(ctx); err != nil {
		slog.Warn("readiness check failed", "error", err)
		apierror.Write(w, http.StatusServiceUnavailable, CodeNotReady)
		return
	}
	jsonhttp.Write(w, http.StatusOK, Status{Status: StatusOK})
}
