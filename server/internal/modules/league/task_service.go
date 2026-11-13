package league

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/module"
)

func (s Service) CreateTask(ctx context.Context, p auth.Principal, leagueID string, in TaskInput) (Task, error) {
	in, err := in.normalized()
	if err != nil {
		return Task{}, err
	}
	var t Task
	err = s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		r, _, err := s.open(ctx, tx, p, leagueID, RoleLead)
		if err != nil {
			return err
		}
		if r.Archived {
			return ErrLeagueArchived
		}
		if t, err = s.tasks.Insert(ctx, tx, leagueID, p.UserID, in); err != nil {
			return err
		}
		if err := s.record(ctx, tx, p, ActionTask, leagueID, map[string]any{"taskId": t.ID, "title": t.Title, "partner": t.Partner}); err != nil {
			return err
		}
		if t.ClosedAt != nil || len(r.Members) == 0 {
			return nil
		}
		return s.notifier.Notify(ctx, tx, p.TenantID, module.Notification{
			Type: NotifyTask, UserIDs: r.Members, Params: map[string]any{"league": r.Title, "task": t.Title}, Link: taskLink(t.ID),
		})
	})
	return t, err
}

func (s Service) UpdateTask(ctx context.Context, p auth.Principal, id string, in TaskInput) (Task, error) {
	in, err := in.normalized()
	if err != nil {
		return Task{}, err
	}
	var t Task
	err = s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		old, err := s.tasks.ByID(ctx, tx, id)
		if err != nil {
			return err
		}
		if _, _, err := s.open(ctx, tx, p, old.LeagueID, RoleLead); err != nil {
			return ErrTaskNotFound
		}
		if t, err = s.tasks.Update(ctx, tx, id, in); err != nil {
			return err
		}
		return s.record(ctx, tx, p, ActionTask, old.LeagueID, map[string]any{"taskId": t.ID, "title": t.Title, "closed": in.Closed})
	})
	return t, err
}

func (s Service) Task(ctx context.Context, p auth.Principal, id string) (TaskDetails, error) {
	var d TaskDetails
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		t, err := s.tasks.ByID(ctx, tx, id)
		if err != nil {
			return err
		}
		l, err := s.leagues.ByID(ctx, tx, t.LeagueID)
		if err != nil {
			return err
		}
		r, _, role, err := s.enter(ctx, tx, p, l.OriginKind, l.OriginID)
		if err != nil {
			return ErrTaskNotFound
		}
		list, err := s.solutions.ByTask(ctx, tx, id)
		if err != nil {
			return err
		}
		d = TaskDetails{Task: t, League: r.Title, Role: role, Solutions: []SolutionRow{}}
		for _, sol := range list {
			if role == RoleStudent && !sol.has(p.UserID) {
				continue
			}
			row := SolutionRow{Solution: sol}
			if role == RoleLead && sol.SimilarTo != nil {
				similar, err := s.solutions.Similar(ctx, tx, *sol.SimilarTo)
				if err != nil {
					return err
				}
				row.Similar = &similar
			}
			d.Solutions = append(d.Solutions, row)
		}
		return nil
	})
	return d, err
}

func (s Service) Explain(ctx context.Context, p auth.Principal, id string, in ExplanationInput) (Solution, error) {
	in, err := in.normalized()
	if err != nil {
		return Solution{}, err
	}
	var out Solution
	err = s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		sol, err := s.solutions.ByID(ctx, tx, id)
		if err != nil {
			return err
		}
		if !sol.has(p.UserID) {
			return ErrSolutionNotFound
		}
		out, err = s.solutions.Explain(ctx, tx, id, in.Text)
		return err
	})
	return out, err
}

func taskLink(id string) string {
	return "/league/tasks/" + id
}
