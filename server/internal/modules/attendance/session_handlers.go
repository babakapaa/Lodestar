package attendance

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type OpenHandler struct{ service Service }

func (h OpenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s, err := h.service.Open(r.Context(), principal(r), r.PathValue("series"), r.PathValue("date"))
	respond(w, s, err)
}

type SessionHandler struct{ service Service }

func (h SessionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s, err := h.service.Session(r.Context(), principal(r), r.PathValue("id"))
	respond(w, s, err)
}

type QRHandler struct{ service Service }

func (h QRHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t, err := h.service.QR(r.Context(), principal(r), r.PathValue("id"))
	respond(w, t, err)
}

type CodeHandler struct{ service Service }

func (h CodeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s, err := h.service.RefreshCode(r.Context(), principal(r), r.PathValue("id"))
	respond(w, s, err)
}

type MarkHandler struct{ service Service }

func (h MarkHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in MarkInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	s, err := h.service.Mark(r.Context(), principal(r), r.PathValue("id"), r.PathValue("user"), in)
	respond(w, s, err)
}

type CloseHandler struct{ service Service }

func (h CloseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s, err := h.service.Close(r.Context(), principal(r), r.PathValue("id"))
	respond(w, s, err)
}
