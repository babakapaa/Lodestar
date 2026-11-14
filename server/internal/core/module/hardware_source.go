package module

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type HardwareSource interface {
	Hardware(ctx context.Context, tx pgx.Tx, tenantID string) ([]HardwareCheck, error)
}
