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
	var in JoinInput
	if !readOptional(w, r, &in) {
		return
	}
	t, err := h.service.Join(r.Context(), principal(r), r.PathValue("id"), in)
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

type CallHookHandler struct {
	service    CallService
	recordings RecordingService
}

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
	room := e.Room.Name
	if e.EgressInfo != nil && room == "" {
		room = e.EgressInfo.RoomName
	}
	tenantID, _, ok := strings.Cut(room, "_")
	if !ok || !validID(tenantID) {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if e.Event == livekit.EgressEnded && e.EgressInfo != nil {
		err = h.recordings.Finish(r.Context(), tenantID, *e.EgressInfo)
	} else {
		err = h.service.Hook(r.Context(), tenantID, e)
	}
	if err != nil {
		slog.Error("call hook", "error", err)
		apierror.Write(w, http.StatusInternalServerError, apierror.CodeInternal)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type StartRecordingHandler struct{ service RecordingService }

func (h StartRecordingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in RecordingInput
	if !readOptional(w, r, &in) {
		return
	}
	c, err := h.service.Start(r.Context(), principal(r), r.PathValue("id"), in)
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, c)
}

type StopRecordingHandler struct{ service RecordingService }

func (h StopRecordingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := h.service.Stop(r.Context(), principal(r), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, c)
}

type ConsentHandler struct{ service RecordingService }

func (h ConsentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in ConsentInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	c, err := h.service.Consent(r.Context(), principal(r), r.PathValue("id"), in)
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, c)
}
