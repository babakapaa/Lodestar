package module

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
)

type Briefer interface {
	Kind() string
	Brief(ctx context.Context, tx pgx.Tx, p auth.Principal, id string) (Brief, error)
}
