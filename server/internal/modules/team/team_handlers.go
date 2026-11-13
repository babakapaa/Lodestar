package team

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type TeamHandler struct{ service Service }

func (h TeamHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d, err := h.service.Get(r.Context(), principal(r), r.PathValue("id"))
	respond(w, http.StatusOK, d, err)
}

type RenameHandler struct{ service Service }

func (h RenameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in NameInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	d, err := h.service.Rename(r.Context(), principal(r), r.PathValue("id"), in)
	respond(w, http.StatusOK, d, err)
}

type LabelHandler struct{ service Service }

func (h LabelHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in LabelInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	d, err := h.service.Label(r.Context(), principal(r), r.PathValue("id"), in)
	respond(w, http.StatusOK, d, err)
}

type InviteHandler struct{ service Service }

func (h InviteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in InviteInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	d, err := h.service.Invite(r.Context(), principal(r), r.PathValue("id"), in)
	respond(w, http.StatusOK, d, err)
}

type UninviteHandler struct{ service Service }

func (h UninviteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Uninvite(r.Context(), principal(r), r.PathValue("id"), r.PathValue("userId")); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type AcceptHandler struct{ service Service }

func (h AcceptHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d, err := h.service.Accept(r.Context(), principal(r), r.PathValue("id"))
	respond(w, http.StatusOK, d, err)
}

type LeaveHandler struct{ service Service }

func (h LeaveHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in LeaveInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	res, err := h.service.Leave(r.Context(), principal(r), r.PathValue("id"), in)
	respond(w, http.StatusOK, res, err)
}

type AddMemberHandler struct{ service Service }

func (h AddMemberHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d, err := h.service.Add(r.Context(), principal(r), r.PathValue("id"), r.PathValue("userId"))
	respond(w, http.StatusOK, d, err)
}

type RemoveMemberHandler struct{ service Service }

func (h RemoveMemberHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	in := RemovalInput{Solo: r.URL.Query().Get("solo") == "true"}
	d, err := h.service.Remove(r.Context(), principal(r), r.PathValue("id"), r.PathValue("userId"), in)
	respond(w, http.StatusOK, d, err)
}

type PreviewHandler struct{ service Service }

func (h PreviewHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	card, err := h.service.Preview(r.Context(), principal(r), r.PathValue("code"))
	respond(w, http.StatusOK, card, err)
}

type JoinHandler struct{ service Service }

func (h JoinHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d, err := h.service.Join(r.Context(), principal(r), r.PathValue("code"))
	respond(w, http.StatusOK, d, err)
}

type DecideHandler struct{ service Service }

func (h DecideHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in DecisionInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	d, err := h.service.Decide(r.Context(), principal(r), r.PathValue("id"), in)
	respond(w, http.StatusOK, d, err)
}
