package structure

import (
	"net/http"
	"strconv"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type HistoryHandler struct {
	service Service
}

func (h HistoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var before int64
	if value := r.URL.Query().Get("before"); value != "" {
		parsed, err := strconv.ParseInt(value, 10, 64)
		if err != nil || parsed < 0 {
			badRequest(w)
			return
		}
		before = parsed
	}
	list, err := h.service.History(r.Context(), principal(r), r.PathValue("id"), before)
	if err != nil {
		writeError(w, err)
		return
	}
	if list == nil {
		list = []Event{}
	}
	jsonhttp.Write(w, http.StatusOK, list)
}
