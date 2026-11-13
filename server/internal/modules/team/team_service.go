package team

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/module"
)

func (s Service) details(ctx context.Context, tx pgx.Tx, p auth.Principal, t Team, ts module.TeamScope, role ScopeRole) (TeamDetails, error) {
	d := TeamDetails{Team: t, Scope: view(ts), Manage: role == RoleLead}
	member := t.has(p.UserID)
	if t.ArchivedAt == nil && !ts.Archived && (role == RoleLead || (member && ts.SelfFormed)) {
		code := t.inviteCode
		d.InviteCode = &code
	}
	invited, err := s.invitations.Invited(ctx, tx, t.ID)
	if err != nil {
		return d, err
	}
	if d.Invited, err = s.people.ByIDs(ctx, tx, invited); err != nil {
		return d, err
	}
	if d.History, err = s.history.ByTeam(ctx, tx, t.ID); err != nil {
		return d, err
	}
	if member {
		if d.Request, err = s.leaves.PendingOf(ctx, tx, t.ID, p.UserID); err != nil {
			return d, err
		}
		d.CanLeaveFreely = freeLeave(ts, time.Now())
	}
	return d, nil
}

func (s Service) Get(ctx context.Context, p auth.Principal, id string) (TeamDetails, error) {
	var d TeamDetails
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		t, ts, role, err := s.team(ctx, tx, p, id, false)
		if err != nil {
			return err
		}
		d, err = s.details(ctx, tx, p, t, ts, role)
		return err
	})
	return d, err
}

func (s Service) Rename(ctx context.Context, p auth.Principal, id string, in NameInput) (TeamDetails, error) {
	normal, err := TeamInput{Name: in.Name}.normalized()
	if err != nil {
		return TeamDetails{}, err
	}
	return s.change(ctx, p, id, func(tx pgx.Tx, t Team, _ module.TeamScope, _ ScopeRole) error {
		return s.teams.Rename(ctx, tx, t.ID, normal.Name)
	})
}

func (s Service) Label(ctx context.Context, p auth.Principal, id string, in LabelInput) (TeamDetails, error) {
	in, err := in.normalized()
	if err != nil {
		return TeamDetails{}, err
	}
	return s.change(ctx, p, id, func(tx pgx.Tx, t Team, _ module.TeamScope, _ ScopeRole) error {
		return s.members.SetLabel(ctx, tx, t.ID, p.UserID, in.Label)
	})
}

func (s Service) Invite(ctx context.Context, p auth.Principal, id string, in InviteInput) (TeamDetails, error) {
	return s.change(ctx, p, id, func(tx pgx.Tx, t Team, ts module.TeamScope, role ScopeRole) error {
		if role == RoleStudent && !ts.SelfFormed {
			return ErrNotSelfFormed
		}
		if ts.Archived || t.ArchivedAt != nil {
			return ErrScopeArchived
		}
		if in.UserID == p.UserID || !slices.Contains(ts.Members, in.UserID) {
			return ErrInvalidTeam
		}
		if len(t.Members) >= ts.MaxSize {
			return ErrTeamFull
		}
		if _, err := s.teams.Of(ctx, tx, ts.Kind, ts.ID, in.UserID); err == nil {
			return ErrInTeam
		} else if !notFound(err) {
			return err
		}
		solo, err := s.separated(ctx, tx, ts, in.UserID)
		if err != nil {
			return err
		}
		if solo {
			return ErrSeparated
		}
		if err := s.invitations.Put(ctx, tx, t.ID, in.UserID, p.UserID); err != nil {
			return err
		}
		inviter, err := s.name(ctx, tx, p.UserID)
		if err != nil {
			return err
		}
		return s.notifier.Notify(ctx, tx, p.TenantID, module.Notification{
			Type: NotifyInvited, UserIDs: []string{in.UserID},
			Params: map[string]any{"scope": ts.Title, "team": t.Name, "inviter": inviter},
			Link:   invitationsLink,
		})
	})
}

func (s Service) Uninvite(ctx context.Context, p auth.Principal, id, userID string) error {
	return s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if userID != p.UserID {
			if _, _, _, err := s.team(ctx, tx, p, id, false); err != nil {
				return err
			}
		}
		removed, err := s.invitations.Delete(ctx, tx, id, userID)
		if err != nil {
			return err
		}
		if !removed {
			return ErrNotInvited
		}
		return nil
	})
}

func (s Service) Accept(ctx context.Context, p auth.Principal, id string) (TeamDetails, error) {
	var d TeamDetails
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		t, err := s.teams.ByID(ctx, tx, id, true)
		if err != nil {
			return err
		}
		ts, role, err := s.scope(ctx, tx, p, t.ScopeKind, t.ScopeID)
		if errors.Is(err, ErrScopeNotFound) || role == RoleLead {
			return ErrNotInvited
		}
		if err != nil {
			return err
		}
		invited, err := s.invitations.Invited(ctx, tx, t.ID)
		if err != nil {
			return err
		}
		if !slices.Contains(invited, p.UserID) {
			return ErrNotInvited
		}
		if err := s.join(ctx, tx, ts, t, p.UserID, EventJoined, p.UserID); err != nil {
			return err
		}
		if t, err = s.teams.ByID(ctx, tx, id, false); err != nil {
			return err
		}
		d, err = s.details(ctx, tx, p, t, ts, role)
		return err
	})
	return d, err
}

func (s Service) Preview(ctx context.Context, p auth.Principal, code string) (TeamCard, error) {
	var card TeamCard
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		t, ts, err := s.byCode(ctx, tx, p, code, false)
		card = TeamCard{Team: t, Scope: view(ts)}
		return err
	})
	return card, err
}

func (s Service) Join(ctx context.Context, p auth.Principal, code string) (TeamDetails, error) {
	var d TeamDetails
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		t, ts, err := s.byCode(ctx, tx, p, code, true)
		if err != nil {
			return err
		}
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
		if err := s.join(ctx, tx, ts, t, p.UserID, EventJoined, p.UserID); err != nil {
			return err
		}
		if t, err = s.teams.ByID(ctx, tx, t.ID, false); err != nil {
			return err
		}
		d, err = s.details(ctx, tx, p, t, ts, RoleStudent)
		return err
	})
	return d, err
}

func (s Service) byCode(ctx context.Context, tx pgx.Tx, p auth.Principal, code string, lock bool) (Team, module.TeamScope, error) {
	t, err := s.teams.ByCode(ctx, tx, code)
	if err != nil {
		return t, module.TeamScope{}, err
	}
	if lock {
		if t, err = s.teams.ByID(ctx, tx, t.ID, true); err != nil {
			return t, module.TeamScope{}, err
		}
	}
	ts, role, err := s.scope(ctx, tx, p, t.ScopeKind, t.ScopeID)
	if errors.Is(err, ErrScopeNotFound) || (err == nil && role != RoleStudent) {
		return t, ts, ErrTeamNotFound
	}
	return t, ts, err
}

func (s Service) Invitations(ctx context.Context, p auth.Principal) ([]TeamCard, error) {
	list := []TeamCard{}
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		rows, err := s.invitations.Mine(ctx, tx, p.UserID)
		if err != nil {
			return err
		}
		for _, row := range rows {
			t, err := s.teams.ByID(ctx, tx, row.teamID, false)
			if err != nil {
				return err
			}
			ts, role, err := s.scope(ctx, tx, p, t.ScopeKind, t.ScopeID)
			if errors.Is(err, ErrScopeNotFound) || role != RoleStudent {
				continue
			}
			if err != nil {
				return err
			}
			card := TeamCard{Team: t, Scope: view(ts)}
			if row.invitedBy != "" {
				people, err := s.people.ByIDs(ctx, tx, []string{row.invitedBy})
				if err != nil {
					return err
				}
				if len(people) == 1 {
					card.InvitedBy = &people[0]
				}
			}
			list = append(list, card)
		}
		return nil
	})
	return list, err
}

func (s Service) Mine(ctx context.Context, p auth.Principal) ([]TeamCard, error) {
	list := []TeamCard{}
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		teams, err := s.teams.ByUser(ctx, tx, p.UserID)
		if err != nil {
			return err
		}
		scopes := map[string]Scope{}
		for _, t := range teams {
			key := t.ScopeKind + "/" + t.ScopeID
			sc, ok := scopes[key]
			if !ok {
				ts, err := s.scopes.Load(ctx, tx, t.ScopeKind, t.ScopeID)
				switch {
				case err == nil:
					sc = view(ts)
				case errors.Is(err, module.ErrTeamScopeNotFound):
					sc = Scope{Kind: t.ScopeKind, ID: t.ScopeID, Archived: true}
				default:
					return err
				}
				scopes[key] = sc
			}
			list = append(list, TeamCard{Team: t, Scope: sc})
		}
		return nil
	})
	return list, err
}

func (s Service) change(ctx context.Context, p auth.Principal, id string, apply func(tx pgx.Tx, t Team, ts module.TeamScope, role ScopeRole) error) (TeamDetails, error) {
	var d TeamDetails
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		t, ts, role, err := s.team(ctx, tx, p, id, true)
		if err != nil {
			return err
		}
		if err := apply(tx, t, ts, role); err != nil {
			return err
		}
		if t, err = s.teams.ByID(ctx, tx, id, false); err != nil {
			return err
		}
		d, err = s.details(ctx, tx, p, t, ts, role)
		return err
	})
	return d, err
}
