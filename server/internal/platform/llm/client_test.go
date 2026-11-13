package llm_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/babakapaa/lodestar/server/internal/platform/llm"
)

func TestCompleteAndStream(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" || r.Header.Get("Authorization") != "Bearer k" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var body struct {
			Model          string `json:"model"`
			Stream         bool   `json:"stream"`
			ResponseFormat *struct {
				Type string `json:"type"`
			} `json:"response_format"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body.Stream {
			w.Header().Set("Content-Type", "text/event-stream")
			_, _ = w.Write([]byte("data: {\"choices\":[{\"delta\":{\"content\":\"При\"}}]}\n\n" +
				"data: {\"choices\":[{\"delta\":{\"content\":\"вет\"}}]}\n\ndata: [DONE]\n\n"))
			return
		}
		content := "ok:" + body.Model
		if body.ResponseFormat != nil {
			content = `{"json":true}`
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"choices": []any{map[string]any{"message": map[string]any{"role": "assistant", "content": content}}}})
	}))
	defer srv.Close()
	c := llm.NewClient(llm.Config{URL: srv.URL + "/v1/", Model: "qwen", Key: "k", Timeout: time.Second})
	ctx := context.Background()
	msgs := []llm.Message{{Role: llm.RoleUser, Content: "привет"}}

	if got, err := c.Complete(ctx, llm.Request{Messages: msgs}); err != nil || got != "ok:qwen" {
		t.Fatalf("complete = %q, %v", got, err)
	}
	if got, err := c.Complete(ctx, llm.Request{Messages: msgs, JSON: true}); err != nil || got != `{"json":true}` {
		t.Fatalf("json = %q, %v", got, err)
	}
	var pieces []string
	got, err := c.Stream(ctx, llm.Request{Messages: msgs}, func(s string) error {
		pieces = append(pieces, s)
		return nil
	})
	if err != nil || got != "Привет" || len(pieces) != 2 {
		t.Fatalf("stream = %q %v, %v", got, pieces, err)
	}

	bad := llm.NewClient(llm.Config{URL: srv.URL + "/v1", Model: "qwen", Timeout: time.Second})
	if _, err := bad.Complete(ctx, llm.Request{Messages: msgs}); !errors.Is(err, llm.ErrModelFailed) {
		t.Fatalf("unauthorized = %v", err)
	}
}
