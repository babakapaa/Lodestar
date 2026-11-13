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

func TestEmbedAndRerank(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Model     string   `json:"model"`
			Input     []string `json:"input"`
			Query     string   `json:"query"`
			Documents []string `json:"documents"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		switch r.URL.Path {
		case "/embeddings":
			data := make([]any, 0, len(body.Input))
			for i := len(body.Input) - 1; i >= 0; i-- {
				data = append(data, map[string]any{"index": i, "embedding": []float32{float32(len(body.Input[i])), 1}})
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"data": data})
		case "/rerank":
			if body.Query == "broken" {
				_ = json.NewEncoder(w).Encode(map[string]any{"results": []any{map[string]any{"index": 9, "relevance_score": 1}}})
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"results": []any{
				map[string]any{"index": 1, "relevance_score": 0.9},
				map[string]any{"index": 0, "relevance_score": 0.1},
			}})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()
	ctx := context.Background()
	cfg := llm.Config{URL: srv.URL, Model: "bge-m3", Timeout: time.Second}

	e := llm.NewEmbedClient(cfg)
	vectors, err := e.Embed(ctx, []string{"a", "bbb"})
	if err != nil || e.Name() != "bge-m3" || len(vectors) != 2 || vectors[0][0] != 1 || vectors[1][0] != 3 {
		t.Fatalf("embed = %v, %v", vectors, err)
	}
	r := llm.NewRerankClient(cfg)
	scores, err := r.Rerank(ctx, "клетка", []string{"x", "y"})
	if err != nil || scores[0] != 0.1 || scores[1] != 0.9 {
		t.Fatalf("rerank = %v, %v", scores, err)
	}
	if _, err := r.Rerank(ctx, "broken", []string{"x"}); !errors.Is(err, llm.ErrModelFailed) {
		t.Fatalf("broken rerank = %v", err)
	}
}
