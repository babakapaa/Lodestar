package navigator

import (
	"net/http"
	"strconv"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type TopicsHandler struct{ service TopicService }

func (h TopicsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	weeks := 0
	if raw := r.URL.Query().Get("weeks"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil {
			badRequest(w)
			return
		}
		weeks = n
	}
	topics, err := h.service.Course(r.Context(), principal(r), r.PathValue("id"), weeks)
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, topics)
}
