package league

import (
	"context"
	"slices"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
)

func (s Service) Decide(ctx context.Context, p auth.Principal, id string, in DecisionInput) (Transition, error) {
	var out Transition
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		t, leagueID, err := s.transitions.ByID(ctx, tx, id)
		if err != nil {
			return err
		}
		r, source, err := s.open(ctx, tx, p, leagueID, RoleLead)
		if err != nil {
			return ErrTransitionNotFound
		}
		at, err := s.sync(ctx, tx, p.TenantID, r, source)
		if err != nil {
			return err
		}
		if t.Status != TransitionPending {
			return ErrTransitionNotFound
		}
		status := TransitionRejected
		if in.Approve {
			status = TransitionApplied
		}
		if out, err = s.transitions.Decide(ctx, tx, t.ID, status, p.UserID); err != nil {
			return err
		}
		if err := s.record(ctx, tx, p, ActionTransition, leagueID, map[string]any{"userId": t.UserID, "from": t.From, "to": t.To, "approve": in.Approve}); err != nil {
			return err
		}
		if !in.Approve || !slices.Contains(r.Members, t.UserID) {
			return nil
		}
		if err := s.players.SetDivision(ctx, tx, leagueID, t.UserID, t.To); err != nil {
			return err
		}
		return s.tellMoved(ctx, tx, p.TenantID, at, out)
	})
	return out, err
}

func (s Service) Adjust(ctx context.Context, p auth.Principal, leagueID string, in AdjustmentInput) (PlayerDetails, error) {
	in, err := in.normalized()
	if err != nil {
		return PlayerDetails{}, err
	}
	var d PlayerDetails
	err = s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		r, source, err := s.open(ctx, tx, p, leagueID, RoleLead)
		if err != nil {
			return err
		}
		if r.Archived {
			return ErrLeagueArchived
		}
		if !slices.Contains(r.Members, in.UserID) {
			return ErrPlayerNotFound
		}
		at, err := s.sync(ctx, tx, p.TenantID, r, source)
		if err != nil {
			return err
		}
		if _, err := s.entries.Adjust(ctx, tx, at.league.ID, at.season.ID, in.UserID, in.Points, in.Reason, p.UserID); err != nil {
			return err
		}
		if err := s.record(ctx, tx, p, ActionAdjusted, at.league.ID, map[string]any{"userId": in.UserID, "points": in.Points, "reason": in.Reason}); err != nil {
			return err
		}
		if err := s.promote(ctx, tx, p.TenantID, at); err != nil {
			return err
		}
		d, err = s.details(ctx, tx, at, in.UserID)
		return err
	})
	return d, err
}

func (s Service) Player(ctx context.Context, p auth.Principal, leagueID, userID string) (PlayerDetails, error) {
	var d PlayerDetails
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		r, source, err := s.open(ctx, tx, p, leagueID, RoleLead)
		if err != nil {
			return err
		}
		if !slices.Contains(r.Members, userID) {
			return ErrPlayerNotFound
		}
		at, err := s.sync(ctx, tx, p.TenantID, r, source)
		if err != nil {
			return err
		}
		d, err = s.details(ctx, tx, at, userID)
		return err
	})
	return d, err
}

func (s Service) details(ctx context.Context, tx pgx.Tx, at synced, userID string) (PlayerDetails, error) {
	b, err := s.board(ctx, tx, at)
	if err != nil {
		return PlayerDetails{}, err
	}
	row, ok := b.find(userID)
	if !ok {
		return PlayerDetails{}, ErrPlayerNotFound
	}
	d := PlayerDetails{
		Person: b.person(userID),
		Standing: Standing{
			Points: row.Points, Division: row.Division, ToNext: at.league.Divisions.toNext(row.Division, row.Points),
			Growth: row.Growth, Display: row.Display, Nickname: row.Nickname, GPUHours: at.league.Divisions[row.Division].GPUHours,
		},
	}
	if d.Grants, err = s.grants.Of(ctx, tx, at.league.ID, &userID); err != nil {
		return d, err
	}
	if d.Entries, err = s.entries.Of(ctx, tx, at.season.ID, userID); err != nil {
		return d, err
	}
	if d.Transitions, err = s.transitions.Of(ctx, tx, at.league.ID, userID); err != nil {
		return d, err
	}
	d.Disputes, err = s.disputes.Of(ctx, tx, at.league.ID, userID)
	return d, err
}
