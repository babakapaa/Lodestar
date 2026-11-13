package team

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type ProposeHandler struct{ service Service }

func (h ProposeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in ProposalInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	d, err := h.service.Propose(r.Context(), principal(r), r.PathValue("id"), in)
	respond(w, http.StatusCreated, d, err)
}

type VoteHandler struct{ service Service }

func (h VoteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in VoteInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	d, err := h.service.Vote(r.Context(), principal(r), r.PathValue("id"), in)
	respond(w, http.StatusOK, d, err)
}

type CardHandler struct{ service Service }

func (h CardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := h.service.Card(r.Context(), principal(r))
	respond(w, http.StatusOK, c, err)
}

type SaveCardHandler struct{ service Service }

func (h SaveCardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in Card
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	c, err := h.service.SaveCard(r.Context(), principal(r), in)
	respond(w, http.StatusOK, c, err)
}

type SuggestHandler struct{ service Service }

func (h SuggestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	out, err := h.service.Suggest(r.Context(), principal(r), r.PathValue("kind"), r.PathValue("id"))
	respond(w, http.StatusOK, out, err)
}
