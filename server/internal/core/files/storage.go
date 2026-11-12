package files

import (
	"context"
	"time"

	"github.com/babakapaa/lodestar/server/internal/platform/objectstore"
)

type Storage interface {
	PresignPut(ctx context.Context, key, contentType string, size int64, ttl time.Duration) (objectstore.Upload, error)
	PresignGet(ctx context.Context, key, name, contentType string, ttl time.Duration) (string, error)
	Size(ctx context.Context, key string) (int64, error)
	Delete(ctx context.Context, key string) error
}
