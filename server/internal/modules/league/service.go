package league

import (
	"context"
	"errors"
	"slices"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/audit"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

type Service struct {
	tenants     database.TenantTx
	rosters     module.Rosters
	scores      module.ScoreSources
	notifier    module.Notifier
	leagues     LeagueStore
	weights     WeightStore
	seasons     SeasonStore
	players     PlayerStore
	entries     EntryStore
	transitions TransitionStore
	disputes    DisputeStore
	people      identity.People
	log         audit.Log
}

func NewService(k module.Kernel) Service {
	return Service{tenants: k.Tenants, rosters: k.Rosters, scores: k.Scores, notifier: k.Notifier}
}

func (s Service) arena(ctx context.Context, tx pgx.Tx, kind, id string) (module.Roster, module.ScoreSource, error) {
	source, ok := s.scores.Find(kind)
	if !ok {
		return module.Roster{}, nil, ErrLeagueNotFound
	}
	rosters, ok := s.rosters.Find(kind)
	if !ok {
		return module.Roster{}, nil, ErrLeagueNotFound
	}
	r, err := rosters.One(ctx, tx, id)
	if errors.Is(err, module.ErrRosterNotFound) {
		return r, nil, ErrLeagueNotFound
	}
	return r, source, err
}

func (s Service) enter(ctx context.Context, tx pgx.Tx, p auth.Principal, kind, id string) (module.Roster, module.ScoreSource, Role, error) {
	r, source, err := s.arena(ctx, tx, kind, id)
	if err != nil {
		return r, nil, "", err
	}
	switch {
	case slices.Contains(r.Leads, p.UserID):
		return r, source, RoleLead, nil
	case slices.Contains(r.Members, p.UserID):
		return r, source, RoleStudent, nil
	}
	return r, nil, "", ErrLeagueNotFound
}

func (s Service) open(ctx context.Context, tx pgx.Tx, p auth.Principal, leagueID string, want Role) (module.Roster, module.ScoreSource, error) {
	l, err := s.leagues.ByID(ctx, tx, leagueID)
	if err != nil {
		return module.Roster{}, nil, err
	}
	r, source, role, err := s.enter(ctx, tx, p, l.OriginKind, l.OriginID)
	if err != nil {
		return r, nil, err
	}
	if role != want {
		return r, nil, ErrLeagueNotFound
	}
	return r, source, nil
}

func (s Service) record(ctx context.Context, tx pgx.Tx, p auth.Principal, action, leagueID string, details map[string]any) error {
	return s.log.Record(ctx, tx, p.UserID, audit.Change{Action: action, TargetKind: "league.league", TargetID: leagueID, Details: details})
}
