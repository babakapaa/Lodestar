package llm_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/platform/llm"
)

func TestPing(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Header.Get("Authorization") != "Bearer k":
			w.WriteHeader(http.StatusUnauthorized)
		case r.URL.Path == "/v1/models":
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()
	ctx := context.Background()
	if err := llm.Ping(ctx, llm.Config{URL: srv.URL + "/v1", Key: "k"}); err != nil {
		t.Fatalf("ok: %v", err)
	}
	if err := llm.Ping(ctx, llm.Config{URL: srv.URL + "/rerank", Key: "k"}); err != nil {
		t.Fatalf("no models list is still reachable: %v", err)
	}
	if err := llm.Ping(ctx, llm.Config{URL: srv.URL + "/v1", Key: "bad"}); !errors.Is(err, llm.ErrModelFailed) {
		t.Fatalf("bad key: %v", err)
	}
}
