package module

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type ReadinessSource interface {
	Problems(ctx context.Context, tx pgx.Tx, tenantID string) ([]Problem, error)
}
