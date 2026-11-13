package team

import (
	"context"
	"slices"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/audit"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/module"
)

func (s Service) View(ctx context.Context, p auth.Principal, kind, id string) (ScopeView, error) {
	var v ScopeView
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		ts, role, err := s.scope(ctx, tx, p, kind, id)
		if err != nil {
			return err
		}
		v, err = s.view(ctx, tx, p, ts, role)
		return err
	})
	return v, err
}

func (s Service) view(ctx context.Context, tx pgx.Tx, p auth.Principal, ts module.TeamScope, role ScopeRole) (ScopeView, error) {
	v := ScopeView{Scope: view(ts), Role: role, Seekers: []Seeker{}}
	mine, err := s.teams.Of(ctx, tx, ts.Kind, ts.ID, p.UserID)
	switch {
	case err == nil:
		v.MyTeam = &mine
	case !notFound(err):
		return v, err
	}
	solos, err := s.solos.ByScope(ctx, tx, ts.Kind, ts.ID)
	if err != nil {
		return v, err
	}
	v.Solo = solos[p.UserID]
	seekers, err := s.seekers.ByScope(ctx, tx, ts.Kind, ts.ID)
	if err != nil {
		return v, err
	}
	v.Seeking = slices.ContainsFunc(seekers, func(x Seeker) bool { return x.UserID == p.UserID })
	if ts.SelfFormed || role == RoleLead {
		v.Seekers = seekers
	}
	if role != RoleLead {
		return v, nil
	}
	if v.Teams, err = s.teams.ByScope(ctx, tx, ts.Kind, ts.ID); err != nil {
		return v, err
	}
	placed, err := s.members.Placed(ctx, tx, ts.Kind, ts.ID)
	if err != nil {
		return v, err
	}
	var unplaced, separated []string
	for _, id := range ts.Members {
		switch {
		case solos[id]:
			separated = append(separated, id)
		case !placed[id]:
			unplaced = append(unplaced, id)
		}
	}
	if v.Unplaced, err = s.people.ByIDs(ctx, tx, unplaced); err != nil {
		return v, err
	}
	if v.Solos, err = s.people.ByIDs(ctx, tx, separated); err != nil {
		return v, err
	}
	v.Requests, err = s.leaves.Pending(ctx, tx, ts.Kind, ts.ID)
	return v, err
}

func (s Service) Create(ctx context.Context, p auth.Principal, kind, id string, in TeamInput) (TeamDetails, error) {
	in, err := in.normalized()
	if err != nil {
		return TeamDetails{}, err
	}
	var d TeamDetails
	err = s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		ts, role, err := s.scope(ctx, tx, p, kind, id)
		if err != nil {
			return err
		}
		if ts.Archived {
			return ErrScopeArchived
		}
		members := in.Members
		if role == RoleStudent {
			if !ts.SelfFormed {
				return ErrNotSelfFormed
			}
			solo, err := s.separated(ctx, tx, ts, p.UserID)
			if err != nil {
				return err
			}
			if solo {
				return ErrSeparated
			}
			members = []string{p.UserID}
		}
		if len(members) > ts.MaxSize {
			return ErrInvalidTeam
		}
		for _, m := range members {
			if !slices.Contains(ts.Members, m) {
				return ErrInvalidTeam
			}
		}
		t, err := s.teams.Insert(ctx, tx, view(ts), in.Name, p.UserID)
		if err != nil {
			return err
		}
		if err := s.history.Add(ctx, tx, t.ID, "", EventCreated, p.UserID); err != nil {
			return err
		}
		event := EventJoined
		if role == RoleLead {
			event = EventAdded
		}
		for _, m := range members {
			if err := s.join(ctx, tx, ts, t, m, event, p.UserID); err != nil {
				return err
			}
			t.Members = append(t.Members, Member{})
		}
		if role == RoleLead {
			if err := s.record(ctx, tx, p, ActionCreated, t.ID, map[string]any{"members": members}); err != nil {
				return err
			}
			if err := s.tellChanged(ctx, tx, p, ts, t, members); err != nil {
				return err
			}
		}
		if t, err = s.teams.ByID(ctx, tx, t.ID, false); err != nil {
			return err
		}
		d, err = s.details(ctx, tx, p, t, ts, role)
		return err
	})
	return d, err
}

func (s Service) Seek(ctx context.Context, p auth.Principal, kind, id string, in SeekingInput) (ScopeView, error) {
	in, err := in.normalized()
	if err != nil {
		return ScopeView{}, err
	}
	var v ScopeView
	err = s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		ts, role, err := s.scope(ctx, tx, p, kind, id)
		if err != nil {
			return err
		}
		if role != RoleStudent {
			return ErrScopeNotFound
		}
		if _, err := s.teams.Of(ctx, tx, kind, id, p.UserID); err == nil {
			return ErrInTeam
		} else if !notFound(err) {
			return err
		}
		if in.Seeking {
			solo, err := s.separated(ctx, tx, ts, p.UserID)
			if err != nil {
				return err
			}
			if solo {
				return ErrSeparated
			}
			err = s.seekers.Put(ctx, tx, kind, id, p.UserID, in.Note)
			if err != nil {
				return err
			}
		} else if err := s.seekers.Delete(ctx, tx, kind, id, p.UserID); err != nil {
			return err
		}
		v, err = s.view(ctx, tx, p, ts, role)
		return err
	})
	return v, err
}

func (s Service) Distribute(ctx context.Context, p auth.Principal, kind, id string) (ScopeView, error) {
	var v ScopeView
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		ts, role, err := s.scope(ctx, tx, p, kind, id)
		if err != nil {
			return err
		}
		if role != RoleLead {
			return ErrScopeNotFound
		}
		if ts.Archived {
			return ErrScopeArchived
		}
		before, err := s.view(ctx, tx, p, ts, role)
		if err != nil {
			return err
		}
		plan := Distribution{MinSize: ts.MinSize, MaxSize: ts.MaxSize}.Plan(before.Teams, before.Unplaced)
		for _, step := range plan {
			t := step.Team
			if t.ID == "" {
				if t, err = s.teams.Insert(ctx, tx, before.Scope, step.Name, p.UserID); err != nil {
					return err
				}
				if err := s.history.Add(ctx, tx, t.ID, "", EventCreated, p.UserID); err != nil {
					return err
				}
			}
			if t, err = s.teams.ByID(ctx, tx, t.ID, true); err != nil {
				return err
			}
			for _, userID := range step.Add {
				if err := s.join(ctx, tx, ts, t, userID, EventAdded, p.UserID); err != nil {
					return err
				}
				t.Members = append(t.Members, Member{})
			}
			if err := s.tellChanged(ctx, tx, p, ts, t, step.Add); err != nil {
				return err
			}
		}
		if len(plan) > 0 {
			if err := s.log.Record(ctx, tx, p.UserID, audit.Change{Action: ActionDistributed, TargetKind: ts.Kind, TargetID: ts.ID, Details: map[string]any{"teams": len(plan)}}); err != nil {
				return err
			}
		}
		v, err = s.view(ctx, tx, p, ts, role)
		return err
	})
	return v, err
}

func (s Service) tellChanged(ctx context.Context, tx pgx.Tx, p auth.Principal, ts module.TeamScope, t Team, userIDs []string) error {
	if len(userIDs) == 0 {
		return nil
	}
	return s.notifier.Notify(ctx, tx, p.TenantID, module.Notification{
		Type: NotifyChanged, UserIDs: userIDs,
		Params: map[string]any{"scope": ts.Title, "team": t.Name},
		Link:   teamLink(t.ID),
	})
}
