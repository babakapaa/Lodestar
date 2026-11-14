package tester

import (
	"log/slog"
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/apierror"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

const (
	CodeUnknownPerson  = "tester_person_unknown"
	CodeInvalidRequest = "invalid_request"
	defaultDevice      = "Режим тестировщика"
	maxDevice          = 100
)

type SignInHandler struct {
	tenantID string
	people   map[string]string
	sessions auth.SessionStore
}

func (h SignInHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req SignInRequest
	if err := jsonhttp.Read(w, r, &req); err != nil {
		apierror.Write(w, http.StatusBadRequest, CodeInvalidRequest)
		return
	}
	userID, ok := h.people[req.Person]
	if !ok {
		apierror.Write(w, http.StatusNotFound, CodeUnknownPerson)
		return
	}
	pair, err := h.sessions.Start(r.Context(), h.tenantID, userID, device(req.Device))
	if err != nil {
		slog.Error("tester sign-in", "error", err)
		apierror.Write(w, http.StatusInternalServerError, apierror.CodeInternal)
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	jsonhttp.Write(w, http.StatusOK, pair)
}

func device(name string) string {
	runes := []rune(name)
	if len(runes) == 0 {
		return defaultDevice
	}
	if len(runes) > maxDevice {
		runes = runes[:maxDevice]
	}
	return string(runes)
}
