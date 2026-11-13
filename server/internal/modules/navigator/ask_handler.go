package navigator

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/babakapaa/lodestar/server/internal/core/apierror"
	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

const answerDeadline = 3 * time.Minute

type AskHandler struct{ service Service }

func (h AskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in MessageInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	_ = http.NewResponseController(w).SetWriteDeadline(time.Now().Add(answerDeadline))
	if !strings.Contains(r.Header.Get("Accept"), "text/event-stream") {
		a, created, err := h.service.Ask(r.Context(), principal(r), r.PathValue("id"), in, nil)
		if err != nil {
			writeError(w, err)
			return
		}
		status := http.StatusOK
		if created {
			status = http.StatusCreated
		}
		jsonhttp.Write(w, status, a)
		return
	}
	stream := newEventStream(w)
	a, _, err := h.service.Ask(r.Context(), principal(r), r.PathValue("id"), in, func(st Stage) error {
		return stream.send("stage", map[string]Stage{"stage": st})
	})
	if err != nil {
		if !stream.started {
			writeError(w, err)
			return
		}
		_ = stream.send("error", apierror.Error{Code: errorCode(err)})
		return
	}
	_ = stream.send("answer", a)
}

func errorCode(err error) string {
	for _, e := range errorCodes {
		if errors.Is(err, e.err) {
			return e.code
		}
	}
	return CodeModelFailed
}
