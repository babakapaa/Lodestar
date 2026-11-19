package module

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type HelpSource interface {
	Help(ctx context.Context, tx pgx.Tx, userID, scopeKind, scopeID string) (HelpSummary, error)
}
