package team

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/module"
)

type Placements struct {
	teams TeamStore
	solos SoloStore
}

func (d Placements) Placement(ctx context.Context, tx pgx.Tx, kind, scopeID, userID string) (module.TeamPlacement, error) {
	var out module.TeamPlacement
	t, err := d.teams.Of(ctx, tx, kind, scopeID, userID)
	switch {
	case err == nil:
		out.TeamID = t.ID
	case !notFound(err):
		return out, err
	}
	solos, err := d.solos.ByScope(ctx, tx, kind, scopeID)
	out.Solo = solos[userID]
	return out, err
}
