package audit

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/babakapaa/lodestar/server/internal/core/apierror"
	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type ListHandler struct {
	service Service
}

func (h ListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	f := Filter{Action: q.Get("action"), ActorID: q.Get("actor")}
	if value := q.Get("before"); value != "" {
		before, err := strconv.ParseInt(value, 10, 64)
		if err != nil || before < 0 {
			apierror.Write(w, http.StatusBadRequest, CodeInvalidRequest)
			return
		}
		f.Before = before
	}
	list, err := h.service.List(r.Context(), principal(r), f)
	if err != nil {
		if errors.Is(err, ErrForbidden) {
			apierror.Write(w, http.StatusForbidden, CodeForbidden)
			return
		}
		slog.Error("read audit log", "error", err)
		apierror.Write(w, http.StatusInternalServerError, apierror.CodeInternal)
		return
	}
	if list == nil {
		list = []Entry{}
	}
	jsonhttp.Write(w, http.StatusOK, list)
}
