package chat

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type ChatsHandler struct{ service Service }

func (h ChatsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.List(r.Context(), principal(r))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, list)
}

type CreateHandler struct{ service Service }

func (h CreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in ChatInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	d, created, err := h.service.Create(r.Context(), principal(r), in)
	writeCreated(w, d, created, err)
}

type DirectHandler struct{ service Service }

func (h DirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in DirectInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	d, created, err := h.service.Direct(r.Context(), principal(r), in)
	writeCreated(w, d, created, err)
}

type ChatHandler struct{ service Service }

func (h ChatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d, err := h.service.Get(r.Context(), principal(r), r.PathValue("id"))
	writeOK(w, d, err)
}

type UpdateHandler struct{ service Service }

func (h UpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in ChatUpdate
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	d, err := h.service.Update(r.Context(), principal(r), r.PathValue("id"), in)
	writeOK(w, d, err)
}

type AddMembersHandler struct{ service Service }

func (h AddMembersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in MembersInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	d, err := h.service.AddMembers(r.Context(), principal(r), r.PathValue("id"), in)
	writeOK(w, d, err)
}

type RemoveMemberHandler struct{ service Service }

func (h RemoveMemberHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.service.RemoveMember(r.Context(), principal(r), r.PathValue("id"), r.PathValue("userId")); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type MuteHandler struct{ service Service }

func (h MuteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in MuteInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	c, err := h.service.Mute(r.Context(), principal(r), r.PathValue("id"), in)
	writeOK(w, c, err)
}

type AcceptHandler struct{ service Service }

func (h AcceptHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := h.service.Accept(r.Context(), principal(r), r.PathValue("id"))
	writeOK(w, c, err)
}

type ReadHandler struct{ service Service }

func (h ReadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in ReadInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	c, err := h.service.Read(r.Context(), principal(r), r.PathValue("id"), in)
	writeOK(w, c, err)
}

func writeOK(w http.ResponseWriter, body any, err error) {
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, body)
}

func writeCreated(w http.ResponseWriter, body any, created bool, err error) {
	if err != nil {
		writeError(w, err)
		return
	}
	status := http.StatusOK
	if created {
		status = http.StatusCreated
	}
	jsonhttp.Write(w, status, body)
}
