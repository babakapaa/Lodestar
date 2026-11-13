package module

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
)

var ErrTeamScopeNotFound = errors.New("team scope not found")

type TeamScopeSource interface {
	Kind() string
	TeamScope(ctx context.Context, tx pgx.Tx, id string) (TeamScope, error)
	PrepareTeamWork(ctx context.Context, tx pgx.Tx, id string, p auth.Principal, work json.RawMessage) (json.RawMessage, error)
	DeliverTeamWork(ctx context.Context, tx pgx.Tx, d TeamDelivery) error
}
