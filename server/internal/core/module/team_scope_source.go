package module

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

var ErrTeamScopeNotFound = errors.New("team scope not found")

type TeamScopeSource interface {
	Kind() string
	TeamScope(ctx context.Context, tx pgx.Tx, id string) (TeamScope, error)
}
