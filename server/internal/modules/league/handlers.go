package league

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type ViewHandler struct{ service Service }

func (h ViewHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v, err := h.service.View(r.Context(), principal(r), r.PathValue("kind"), r.PathValue("id"))
	respond(w, http.StatusOK, v, err)
}

type RulesHandler struct{ service Service }

func (h RulesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in RulesInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	v, err := h.service.SetRules(r.Context(), principal(r), r.PathValue("id"), in)
	respond(w, http.StatusOK, v, err)
}

type DisplayHandler struct{ service Service }

func (h DisplayHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in DisplayInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	v, err := h.service.SetDisplay(r.Context(), principal(r), r.PathValue("id"), in)
	respond(w, http.StatusOK, v, err)
}

type SeasonHandler struct{ service Service }

func (h SeasonHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v, err := h.service.CloseSeason(r.Context(), principal(r), r.PathValue("id"))
	respond(w, http.StatusOK, v, err)
}

type AdjustHandler struct{ service Service }

func (h AdjustHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in AdjustmentInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	d, err := h.service.Adjust(r.Context(), principal(r), r.PathValue("id"), in)
	respond(w, http.StatusCreated, d, err)
}

type PlayerHandler struct{ service Service }

func (h PlayerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d, err := h.service.Player(r.Context(), principal(r), r.PathValue("id"), r.PathValue("userId"))
	respond(w, http.StatusOK, d, err)
}

type DisputeHandler struct{ service Service }

func (h DisputeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in DisputeInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	d, err := h.service.Dispute(r.Context(), principal(r), r.PathValue("id"), in)
	respond(w, http.StatusCreated, d, err)
}

type DecideHandler struct{ service Service }

func (h DecideHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in DecisionInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	t, err := h.service.Decide(r.Context(), principal(r), r.PathValue("id"), in)
	respond(w, http.StatusOK, t, err)
}

type AnswerHandler struct{ service Service }

func (h AnswerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in AnswerInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	d, err := h.service.Answer(r.Context(), principal(r), r.PathValue("id"), in)
	respond(w, http.StatusOK, d, err)
}
