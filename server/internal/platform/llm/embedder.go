package llm

import "context"

type Embedder interface {
	Name() string
	Embed(ctx context.Context, texts []string) ([][]float32, error)
}
