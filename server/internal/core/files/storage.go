package files

import (
	"context"
	"io"
	"time"

	"github.com/babakapaa/lodestar/server/internal/platform/objectstore"
)

type Storage interface {
	PresignPut(ctx context.Context, key, contentType string, size int64, ttl time.Duration) (objectstore.Upload, error)
	PresignGet(ctx context.Context, key, name, contentType string, ttl time.Duration) (string, error)
	Size(ctx context.Context, key string) (int64, error)
	Open(ctx context.Context, key string) (io.ReadCloser, error)
	Put(ctx context.Context, key, contentType string, body io.Reader, size int64) error
	Delete(ctx context.Context, key string) error
}
