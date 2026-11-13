package chat

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type MessagesHandler struct{ service MessageService }

func (h MessagesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q, err := ParseMessageQuery(r.URL.Query())
	if err != nil {
		writeError(w, err)
		return
	}
	page, err := h.service.Page(r.Context(), principal(r), r.PathValue("id"), q)
	writeOK(w, page, err)
}

type SendHandler struct{ service MessageService }

func (h SendHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in MessageInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	m, created, err := h.service.Send(r.Context(), principal(r), r.PathValue("id"), in)
	writeCreated(w, m, created, err)
}

type EditHandler struct{ service MessageService }

func (h EditHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in EditInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	m, err := h.service.Edit(r.Context(), principal(r), r.PathValue("id"), in)
	writeOK(w, m, err)
}

type DeleteMessageHandler struct{ service MessageService }

func (h DeleteMessageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m, err := h.service.Delete(r.Context(), principal(r), r.PathValue("id"))
	writeOK(w, m, err)
}
