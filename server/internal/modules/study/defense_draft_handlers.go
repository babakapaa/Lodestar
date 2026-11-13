package study

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type DraftQuestionsHandler struct{ service DefenseService }

func (h DraftQuestionsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	questions, err := h.service.DraftQuestions(r.Context(), principal(r), r.PathValue("id"), r.PathValue("student"))
	respond(w, http.StatusOK, map[string]any{"questions": questions}, err)
}

type DraftProtocolHandler struct{ service DefenseService }

func (h DraftProtocolHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in DefenseInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	summary, err := h.service.DraftProtocol(r.Context(), principal(r), r.PathValue("id"), r.PathValue("student"), in)
	respond(w, http.StatusOK, map[string]any{"summary": summary}, err)
}
