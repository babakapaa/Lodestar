package search

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
)

type Source interface {
	Kind() string
	Search(ctx context.Context, tx pgx.Tx, p auth.Principal, q Query, limit int) ([]Hit, error)
}
