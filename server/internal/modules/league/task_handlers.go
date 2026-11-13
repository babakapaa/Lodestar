package league

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type CreateTaskHandler struct{ service Service }

func (h CreateTaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in TaskInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	t, err := h.service.CreateTask(r.Context(), principal(r), r.PathValue("id"), in)
	respond(w, http.StatusCreated, t, err)
}

type TaskHandler struct{ service Service }

func (h TaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d, err := h.service.Task(r.Context(), principal(r), r.PathValue("id"))
	respond(w, http.StatusOK, d, err)
}

type UpdateTaskHandler struct{ service Service }

func (h UpdateTaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in TaskInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	t, err := h.service.UpdateTask(r.Context(), principal(r), r.PathValue("id"), in)
	respond(w, http.StatusOK, t, err)
}

type ReviewHandler struct{ service Service }

func (h ReviewHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in ReviewInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	s, err := h.service.Review(r.Context(), principal(r), r.PathValue("id"), in)
	respond(w, http.StatusOK, s, err)
}

type ExplainHandler struct{ service Service }

func (h ExplainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in ExplanationInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	s, err := h.service.Explain(r.Context(), principal(r), r.PathValue("id"), in)
	respond(w, http.StatusOK, s, err)
}
