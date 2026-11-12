package live

import (
	"log/slog"
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/apierror"
	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type EventsHandler struct {
	backlog Backlog
}

func (h EventsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	after, ok := afterParam(r)
	if !ok {
		apierror.Write(w, http.StatusBadRequest, "invalid_request")
		return
	}
	events, more, err := h.backlog.After(r.Context(), principal(r), after)
	if err != nil {
		slog.Error("read live backlog", "error", err)
		apierror.Write(w, http.StatusInternalServerError, apierror.CodeInternal)
		return
	}
	if events == nil {
		events = []Event{}
	}
	jsonhttp.Write(w, http.StatusOK, EventsPage{Events: events, More: more})
}
