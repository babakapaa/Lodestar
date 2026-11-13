package team

import (
	"context"
	"slices"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/module"
)

func (s Service) Add(ctx context.Context, p auth.Principal, id, userID string) (TeamDetails, error) {
	return s.change(ctx, p, id, func(tx pgx.Tx, t Team, ts module.TeamScope, role ScopeRole) error {
		if role != RoleLead {
			return ErrTeamNotFound
		}
		if !slices.Contains(ts.Members, userID) {
			return ErrInvalidTeam
		}
		if t.has(userID) {
			return nil
		}
		event, action := EventAdded, ActionMemberAdded
		from, err := s.teams.Of(ctx, tx, ts.Kind, ts.ID, userID)
		switch {
		case err == nil:
			if from, err = s.teams.ByID(ctx, tx, from.ID, true); err != nil {
				return err
			}
			if err := s.drop(ctx, tx, from, userID, EventMovedOut, p.UserID); err != nil {
				return err
			}
			event, action = EventMovedIn, ActionMemberMoved
		case !notFound(err):
			return err
		}
		if err := s.join(ctx, tx, ts, t, userID, event, p.UserID); err != nil {
			return err
		}
		if err := s.record(ctx, tx, p, action, t.ID, map[string]any{"student": userID}); err != nil {
			return err
		}
		return s.tellChanged(ctx, tx, p, ts, t, []string{userID})
	})
}

func (s Service) Remove(ctx context.Context, p auth.Principal, id, userID string, in RemovalInput) (TeamDetails, error) {
	return s.change(ctx, p, id, func(tx pgx.Tx, t Team, ts module.TeamScope, role ScopeRole) error {
		if role != RoleLead {
			return ErrTeamNotFound
		}
		if !t.has(userID) {
			return ErrInvalidTeam
		}
		event, action := EventRemoved, ActionRemoved
		if in.Solo {
			event, action = EventSeparated, ActionSeparated
			if err := s.solos.Put(ctx, tx, ts.Kind, ts.ID, userID, p.UserID); err != nil {
				return err
			}
		}
		if err := s.drop(ctx, tx, t, userID, event, p.UserID); err != nil {
			return err
		}
		if err := s.record(ctx, tx, p, action, t.ID, map[string]any{"student": userID}); err != nil {
			return err
		}
		return s.tellChanged(ctx, tx, p, ts, t, []string{userID})
	})
}
