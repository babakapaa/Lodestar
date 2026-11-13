package chat

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/babakapaa/lodestar/server/internal/core/apierror"
	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
	"github.com/babakapaa/lodestar/server/internal/platform/livekit"
)

type CallHandler struct{ service CallService }

func (h CallHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := h.service.Active(r.Context(), principal(r), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, c)
}

type JoinCallHandler struct{ service CallService }

func (h JoinCallHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t, err := h.service.Join(r.Context(), principal(r), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, t)
}

type LeaveCallHandler struct{ service CallService }

func (h LeaveCallHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Leave(r.Context(), principal(r), r.PathValue("id")); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type CallHookHandler struct{ service CallService }

func (h CallHookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	issuer := h.service.k.Calls
	if issuer == nil {
		apierror.NotFound(w, r)
		return
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		badRequest(w)
		return
	}
	if err := issuer.Verify(r.Header.Get("Authorization"), body, time.Now()); err != nil {
		apierror.Write(w, http.StatusUnauthorized, CodeBadHook)
		return
	}
	var e livekit.Event
	if err := json.Unmarshal(body, &e); err != nil {
		badRequest(w)
		return
	}
	tenantID, _, ok := strings.Cut(e.Room.Name, "_")
	if !ok || !validID(tenantID) {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err := h.service.Hook(r.Context(), tenantID, e); err != nil {
		slog.Error("call hook", "error", err)
		apierror.Write(w, http.StatusInternalServerError, apierror.CodeInternal)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
