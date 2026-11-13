package team

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/module"
)

const RosterTeam = "team.team"

type TeamRosters struct {
	scopes module.TeamScopes
	teams  TeamStore
}

func (TeamRosters) Kind() string {
	return RosterTeam
}

func (r TeamRosters) All(ctx context.Context, tx pgx.Tx) ([]module.Roster, error) {
	return r.load(ctx, tx, nil)
}

func (r TeamRosters) One(ctx context.Context, tx pgx.Tx, id string) (module.Roster, error) {
	list, err := r.load(ctx, tx, []string{id})
	if err != nil {
		return module.Roster{}, err
	}
	if len(list) == 0 {
		return module.Roster{}, module.ErrRosterNotFound
	}
	return list[0], nil
}

func (r TeamRosters) load(ctx context.Context, tx pgx.Tx, ids []string) ([]module.Roster, error) {
	teams, err := r.teams.All(ctx, tx, ids)
	if err != nil {
		return nil, err
	}
	archived := map[string]bool{}
	list := make([]module.Roster, 0, len(teams))
	for _, t := range teams {
		key := t.ScopeKind + "/" + t.ScopeID
		gone, ok := archived[key]
		if !ok {
			ts, err := r.scopes.Load(ctx, tx, t.ScopeKind, t.ScopeID)
			switch {
			case errors.Is(err, module.ErrTeamScopeNotFound):
				gone = true
			case err != nil:
				return nil, err
			default:
				gone = ts.Archived
			}
			archived[key] = gone
		}
		list = append(list, module.Roster{
			Kind: RosterTeam, ID: t.ID, Title: t.Name, Members: t.ids(), Leads: []string{},
			Group: true, Archived: gone || t.ArchivedAt != nil,
		})
	}
	return list, nil
}
