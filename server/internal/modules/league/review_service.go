package league

import (
	"context"
	"slices"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/module"
)

func (s Service) Review(ctx context.Context, p auth.Principal, id string, in ReviewInput) (Solution, error) {
	var out Solution
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		sol, err := s.solutions.ByID(ctx, tx, id)
		if err != nil {
			return err
		}
		t, err := s.tasks.ByID(ctx, tx, sol.TaskID)
		if err != nil {
			return err
		}
		r, source, err := s.open(ctx, tx, p, t.LeagueID, RoleLead)
		if err != nil {
			return ErrSolutionNotFound
		}
		if r.Archived {
			return ErrLeagueArchived
		}
		if in, err = in.normalized(sol, t); err != nil {
			return err
		}
		at, err := s.sync(ctx, tx, p.TenantID, r, source)
		if err != nil {
			return err
		}
		if err := s.award(ctx, tx, at, t, sol, in); err != nil {
			return err
		}
		if out, err = s.solutions.Review(ctx, tx, id, p.UserID, in); err != nil {
			return err
		}
		if err := s.record(ctx, tx, p, ActionReview, at.league.ID, map[string]any{"solutionId": id, "points": in.Points, "passed": in.Passed}); err != nil {
			return err
		}
		if err := s.promote(ctx, tx, p.TenantID, at); err != nil {
			return err
		}
		return s.notifier.Notify(ctx, tx, p.TenantID, module.Notification{
			Type: NotifyReviewed, UserIDs: sol.Members, Params: map[string]any{"league": r.Title, "task": t.Title}, Link: taskLink(t.ID),
		})
	})
	return out, err
}

func (s Service) award(ctx context.Context, tx pgx.Tx, at synced, t Task, sol Solution, in ReviewInput) error {
	existing, err := s.entries.Works(ctx, tx, at.league.ID, EntryTask)
	if err != nil {
		return err
	}
	kind, maxPoints := TaskScopeKind, float64(t.Points)
	score := float64(in.Points)
	for _, userID := range sol.Members {
		old, found := existing[workKey(userID, TaskScopeKind, t.ID)]
		if found && old.SeasonID != at.season.ID {
			continue
		}
		passed := slices.Contains(in.Passed, userID) && in.Points > 0 && slices.Contains(at.roster.Members, userID)
		if !passed {
			if found {
				if err := s.entries.Delete(ctx, tx, old.ID); err != nil {
					return err
				}
			}
			continue
		}
		next := Entry{
			ID: old.ID, UserID: userID, Kind: EntryTask, WorkKind: &kind, WorkID: &t.ID, Title: t.Title,
			Score: &score, Max: &maxPoints, Weight: &t.Points, Points: in.Points, Reason: in.Comment,
		}
		if found {
			err = s.entries.UpdateWork(ctx, tx, next)
		} else {
			err = s.entries.InsertWork(ctx, tx, at.league.ID, at.season.ID, next)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
