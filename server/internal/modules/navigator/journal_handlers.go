package navigator

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type JournalHandler struct{ service Service }

func (h JournalHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	list, err := h.service.Journal(r.Context(), principal(r), q.Get("scopeKind"), q.Get("scopeId"))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, list)
}

type ForgetHandler struct{ service Service }

func (h ForgetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Forget(r.Context(), principal(r), r.PathValue("id")); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type SupportHandler struct{ service Service }

func (h SupportHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	out, err := h.service.Support(r.Context(), principal(r))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, out)
}
