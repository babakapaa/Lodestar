package module

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type TeamDirectory interface {
	Placement(ctx context.Context, tx pgx.Tx, kind, scopeID, userID string) (TeamPlacement, error)
	Teams(ctx context.Context, tx pgx.Tx, userID string) ([]TeamRef, error)
}
