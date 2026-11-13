package team

import (
	"context"
	"fmt"

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

func (d Placements) Teams(ctx context.Context, tx pgx.Tx, userID string) ([]module.TeamRef, error) {
	rows, err := tx.Query(ctx, `
		SELECT t.id::text, t.name, array_agg(m2.user_id::text ORDER BY m2.joined_at)
		FROM team.members m
		JOIN team.teams t ON t.id = m.team_id AND t.archived_at IS NULL
		JOIN team.members m2 ON m2.team_id = t.id
		WHERE m.user_id::text = $1
		GROUP BY t.id, t.name, t.created_at ORDER BY t.created_at`, userID)
	if err != nil {
		return nil, fmt.Errorf("teams of user: %w", err)
	}
	defer rows.Close()
	var out []module.TeamRef
	for rows.Next() {
		var t module.TeamRef
		if err := rows.Scan(&t.ID, &t.Name, &t.Members); err != nil {
			return nil, fmt.Errorf("teams of user: %w", err)
		}
		out = append(out, t)
	}
	return out, rows.Err()
}
