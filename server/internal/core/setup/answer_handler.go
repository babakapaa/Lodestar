package setup

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type AnswerHandler struct {
	service Service
}

func (h AnswerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in AnswerInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	a, err := h.service.Answer(r.Context(), principal(r), r.PathValue("question"), in)
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, a)
}
