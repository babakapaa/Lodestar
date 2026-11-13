package llm

import "context"

type Model interface {
	Complete(ctx context.Context, req Request) (string, error)
	Stream(ctx context.Context, req Request, delta func(string) error) (string, error)
}
