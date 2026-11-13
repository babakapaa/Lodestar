package chat_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
)

type fakeEgress struct {
	server  *httptest.Server
	mu      sync.Mutex
	started []map[string]any
	stopped []string
}

func newFakeEgress(t *testing.T) *fakeEgress {
	f := &fakeEgress{}
	f.server = httptest.NewServer(http.HandlerFunc(f.serve))
	t.Cleanup(f.server.Close)
	return f
}

func (f *fakeEgress) serve(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var in map[string]any
	_ = json.Unmarshal(body, &in)
	f.mu.Lock()
	defer f.mu.Unlock()
	switch {
	case strings.HasSuffix(r.URL.Path, "/StartRoomCompositeEgress"):
		f.started = append(f.started, in)
		_, _ = w.Write([]byte(`{"egressId":"EG_` + strconv.Itoa(len(f.started)) + `"}`))
	case strings.HasSuffix(r.URL.Path, "/StopEgress"):
		f.stopped = append(f.stopped, in["egress_id"].(string))
		_, _ = w.Write([]byte(`{}`))
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (f *fakeEgress) key(i int) string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.started[i]["file_outputs"].([]any)[0].(map[string]any)["filepath"].(string)
}
