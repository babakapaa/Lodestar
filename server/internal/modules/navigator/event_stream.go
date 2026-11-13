package navigator

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type eventStream struct {
	w       http.ResponseWriter
	rc      *http.ResponseController
	started bool
}

func newEventStream(w http.ResponseWriter) *eventStream {
	return &eventStream{w: w, rc: http.NewResponseController(w)}
}

func (s *eventStream) send(event string, data any) error {
	if !s.started {
		s.w.Header().Set("Content-Type", "text/event-stream")
		s.w.Header().Set("Cache-Control", "no-store")
		s.w.Header().Set("X-Accel-Buffering", "no")
		s.w.WriteHeader(http.StatusOK)
		s.started = true
	}
	raw, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(s.w, "event: %s\ndata: %s\n\n", event, raw); err != nil {
		return err
	}
	return s.rc.Flush()
}
