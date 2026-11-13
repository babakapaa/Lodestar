package league

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
)

func (s Service) SetRules(ctx context.Context, p auth.Principal, leagueID string, in RulesInput) (LeagueView, error) {
	var v LeagueView
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		r, source, err := s.open(ctx, tx, p, leagueID, RoleLead)
		if err != nil {
			return err
		}
		at, err := s.sync(ctx, tx, p.TenantID, r, source)
		if err != nil {
			return err
		}
		if r.Archived {
			return ErrLeagueArchived
		}
		weights, err := s.weights.All(ctx, tx, at.league.ID)
		if err != nil {
			return err
		}
		gpuCap, err := s.gpuCap(ctx, tx)
		if err != nil {
			return err
		}
		if in, err = in.normalized(weights.formula(at.works), gpuCap); err != nil {
			return err
		}
		if err := s.leagues.SetRules(ctx, tx, at.league.ID, in.Divisions, in.Confirm); err != nil {
			return err
		}
		if err := s.weights.Replace(ctx, tx, at.league.ID, in.Weights); err != nil {
			return err
		}
		if err := s.players.Clamp(ctx, tx, at.league.ID, len(in.Divisions)-1); err != nil {
			return err
		}
		if err := s.record(ctx, tx, p, ActionRules, at.league.ID, map[string]any{"divisions": in.Divisions, "confirm": in.Confirm, "weights": in.Weights}); err != nil {
			return err
		}
		v, err = s.reload(ctx, tx, p, at)
		return err
	})
	return v, err
}

func (s Service) SetDisplay(ctx context.Context, p auth.Principal, leagueID string, in DisplayInput) (LeagueView, error) {
	in, err := in.normalized()
	if err != nil {
		return LeagueView{}, err
	}
	var v LeagueView
	err = s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		r, source, err := s.open(ctx, tx, p, leagueID, RoleStudent)
		if err != nil {
			return err
		}
		at, err := s.sync(ctx, tx, p.TenantID, r, source)
		if err != nil {
			return err
		}
		if err := s.players.Ensure(ctx, tx, at.league.ID, []string{p.UserID}); err != nil {
			return err
		}
		if err := s.players.SetDisplay(ctx, tx, at.league.ID, p.UserID, in); err != nil {
			return err
		}
		v, err = s.view(ctx, tx, p, RoleStudent, at)
		return err
	})
	return v, err
}
