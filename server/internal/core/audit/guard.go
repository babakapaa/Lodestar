package audit

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
)

type Guard interface {
	AllowedInTenant(ctx context.Context, tx pgx.Tx, p auth.Principal, permission string) (bool, error)
}
