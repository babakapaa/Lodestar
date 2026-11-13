package module

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
)

var ErrLinkTargetNotFound = errors.New("link target not found")

type LinkSource interface {
	Kind() string
	Linkable(ctx context.Context, tx pgx.Tx, p auth.Principal) ([]LinkTarget, error)
	Describe(ctx context.Context, tx pgx.Tx, id string) (LinkTarget, error)
	CanLink(ctx context.Context, tx pgx.Tx, p auth.Principal, id string) error
}
