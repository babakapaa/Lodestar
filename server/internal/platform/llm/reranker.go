package llm

import "context"

type Reranker interface {
	Rerank(ctx context.Context, query string, documents []string) ([]float64, error)
}
