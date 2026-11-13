package league

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type AwardsHandler struct{ service Service }

func (h AwardsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.Awards(r.Context(), principal(r))
	respond(w, http.StatusOK, list, err)
}

type CreateAwardHandler struct{ service Service }

func (h CreateAwardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in AwardInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	a, err := h.service.CreateAward(r.Context(), principal(r), in)
	respond(w, http.StatusCreated, a, err)
}

type UpdateAwardHandler struct{ service Service }

func (h UpdateAwardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in AwardInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	a, err := h.service.UpdateAward(r.Context(), principal(r), r.PathValue("id"), in)
	respond(w, http.StatusOK, a, err)
}

type GrantHandler struct{ service Service }

func (h GrantHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in GrantInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	g, err := h.service.Grant(r.Context(), principal(r), r.PathValue("id"), in)
	respond(w, http.StatusCreated, g, err)
}

type RevokeHandler struct{ service Service }

func (h RevokeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Revoke(r.Context(), principal(r), r.PathValue("id")); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
