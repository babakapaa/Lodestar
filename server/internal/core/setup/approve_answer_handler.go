package setup

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type ApproveAnswerHandler struct {
	service Service
}

func (h ApproveAnswerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a, err := h.service.ApproveAnswer(r.Context(), principal(r), r.PathValue("question"))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, a)
}
