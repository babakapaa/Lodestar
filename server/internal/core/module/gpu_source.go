package module

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type GPUSource interface {
	GPUGrants(ctx context.Context, tx pgx.Tx, userID string, from, to time.Time) ([]GPUGrant, error)
}
