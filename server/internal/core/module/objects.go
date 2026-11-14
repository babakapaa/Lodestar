package module

import (
	"context"
	"io"
)

type Objects interface {
	Open(ctx context.Context, key string) (io.ReadCloser, error)
	Put(ctx context.Context, key, contentType string, body io.Reader, size int64) error
}
