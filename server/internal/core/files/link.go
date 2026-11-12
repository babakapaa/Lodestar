package files

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
)

type Link interface {
	Readable(ctx context.Context, tx pgx.Tx, p auth.Principal, fileID string) (bool, error)
	InUse(ctx context.Context, tx pgx.Tx, fileID string) (bool, error)
}
