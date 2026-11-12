package privacy

import (
	"context"
	"io"
)

type Blob struct {
	Path string
	Open func(ctx context.Context) (io.ReadCloser, error)
}
