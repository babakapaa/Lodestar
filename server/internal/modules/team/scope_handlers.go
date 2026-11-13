package team

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type ScopeHandler struct{ service Service }

func (h ScopeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v, err := h.service.View(r.Context(), principal(r), r.PathValue("kind"), r.PathValue("id"))
	respond(w, http.StatusOK, v, err)
}

type CreateHandler struct{ service Service }

func (h CreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in TeamInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	d, err := h.service.Create(r.Context(), principal(r), r.PathValue("kind"), r.PathValue("id"), in)
	respond(w, http.StatusCreated, d, err)
}

type SeekHandler struct{ service Service }

func (h SeekHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in SeekingInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	v, err := h.service.Seek(r.Context(), principal(r), r.PathValue("kind"), r.PathValue("id"), in)
	respond(w, http.StatusOK, v, err)
}

type DistributeHandler struct{ service Service }

func (h DistributeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v, err := h.service.Distribute(r.Context(), principal(r), r.PathValue("kind"), r.PathValue("id"))
	respond(w, http.StatusOK, v, err)
}

type MineHandler struct{ service Service }

func (h MineHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.Mine(r.Context(), principal(r))
	respond(w, http.StatusOK, list, err)
}

type InvitationsHandler struct{ service Service }

func (h InvitationsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.Invitations(r.Context(), principal(r))
	respond(w, http.StatusOK, list, err)
}
