package llm

import (
	"context"
	"io"
)

type Transcriber interface {
	Transcribe(ctx context.Context, audio io.Reader, name string) (string, error)
}
