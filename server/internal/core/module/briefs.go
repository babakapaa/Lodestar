package module

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
)

var ErrBriefNotFound = errors.New("nothing to brief about")

type Briefs []Briefer

func (b Briefs) Brief(ctx context.Context, tx pgx.Tx, p auth.Principal, kind, id string) (Brief, error) {
	for _, briefer := range b {
		if briefer.Kind() == kind {
			return briefer.Brief(ctx, tx, p, id)
		}
	}
	return Brief{}, ErrBriefNotFound
}
