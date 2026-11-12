package privacy

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Section interface {
	Name() string
	Export(ctx context.Context, tx pgx.Tx, userID string) (any, error)
}

type Eraser interface {
	Erase(ctx context.Context, tx pgx.Tx, tenantID, userID string) error
}

type BlobSource interface {
	Blobs(ctx context.Context, tx pgx.Tx, tenantID, userID string) ([]Blob, error)
}
