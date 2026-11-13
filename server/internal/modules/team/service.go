package team

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/audit"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/llm"
)

type Service struct {
	tenants     database.TenantTx
	scopes      module.TeamScopes
	notifier    module.Notifier
	teams       TeamStore
	members     MemberStore
	history     HistoryStore
	invitations InvitationStore
	seekers     SeekerStore
	solos       SoloStore
	leaves      LeaveStore
	proposals   ProposalStore
	cards       CardStore
	model       llm.Model
	people      identity.People
	log         audit.Log
}

func NewService(k module.Kernel) Service {
	return Service{tenants: k.Tenants, scopes: k.TeamScopes, notifier: k.Notifier, model: k.Model}
}

func (s Service) scope(ctx context.Context, tx pgx.Tx, p auth.Principal, kind, id string) (module.TeamScope, ScopeRole, error) {
	ts, err := s.scopes.Load(ctx, tx, kind, id)
	if errors.Is(err, module.ErrTeamScopeNotFound) {
		return ts, "", ErrScopeNotFound
	}
	if err != nil {
		return ts, "", err
	}
	switch {
	case slices.Contains(ts.Leads, p.UserID):
		return ts, RoleLead, nil
	case slices.Contains(ts.Members, p.UserID):
		return ts, RoleStudent, nil
	}
	return ts, "", ErrScopeNotFound
}

func (s Service) team(ctx context.Context, tx pgx.Tx, p auth.Principal, id string, lock bool) (Team, module.TeamScope, ScopeRole, error) {
	t, err := s.teams.ByID(ctx, tx, id, lock)
	if err != nil {
		return t, module.TeamScope{}, "", err
	}
	ts, role, err := s.scope(ctx, tx, p, t.ScopeKind, t.ScopeID)
	if errors.Is(err, ErrScopeNotFound) {
		return t, ts, "", ErrTeamNotFound
	}
	if err != nil {
		return t, ts, "", err
	}
	if role == RoleStudent && !t.has(p.UserID) {
		return t, ts, "", ErrTeamNotFound
	}
	return t, ts, role, nil
}

func (s Service) join(ctx context.Context, tx pgx.Tx, ts module.TeamScope, t Team, userID string, event HistoryEvent, actorID string) error {
	if ts.Archived {
		return ErrScopeArchived
	}
	if t.ArchivedAt != nil {
		return ErrTeamNotFound
	}
	if len(t.Members) >= ts.MaxSize {
		return ErrTeamFull
	}
	if err := s.members.Add(ctx, tx, t, userID); err != nil {
		return err
	}
	if err := s.seekers.Delete(ctx, tx, ts.Kind, ts.ID, userID); err != nil {
		return err
	}
	if err := s.solos.Delete(ctx, tx, ts.Kind, ts.ID, userID); err != nil {
		return err
	}
	if err := s.invitations.DeleteInScope(ctx, tx, ts.Kind, ts.ID, userID); err != nil {
		return err
	}
	return s.history.Add(ctx, tx, t.ID, userID, event, actorID)
}

func (s Service) drop(ctx context.Context, tx pgx.Tx, t Team, userID string, event HistoryEvent, actorID string) error {
	if err := s.members.Remove(ctx, tx, t.ID, userID); err != nil {
		return err
	}
	if err := s.leaves.Close(ctx, tx, t.ID, userID); err != nil {
		return err
	}
	if err := s.history.Add(ctx, tx, t.ID, userID, event, actorID); err != nil {
		return err
	}
	if len(t.Members) <= 1 {
		return s.teams.Archive(ctx, tx, t.ID)
	}
	return nil
}

func (s Service) separated(ctx context.Context, tx pgx.Tx, ts module.TeamScope, userID string) (bool, error) {
	solos, err := s.solos.ByScope(ctx, tx, ts.Kind, ts.ID)
	return solos[userID], err
}

func (s Service) record(ctx context.Context, tx pgx.Tx, p auth.Principal, action, teamID string, details map[string]any) error {
	return s.log.Record(ctx, tx, p.UserID, audit.Change{Action: action, TargetKind: "team.team", TargetID: teamID, Details: details})
}

func (s Service) name(ctx context.Context, tx pgx.Tx, userID string) (string, error) {
	people, err := s.people.ByIDs(ctx, tx, []string{userID})
	if err != nil || len(people) == 0 {
		return "", err
	}
	if people[0].DisplayName != "" {
		return people[0].DisplayName, nil
	}
	return people[0].Email, nil
}

func view(ts module.TeamScope) Scope {
	return Scope{
		Kind: ts.Kind, ID: ts.ID, Title: ts.Title, Context: ts.Context, Link: ts.Link,
		MinSize: ts.MinSize, MaxSize: ts.MaxSize, SelfFormed: ts.SelfFormed,
		FreeLeaveUntil: ts.FreeLeaveUntil, DueAt: ts.DueAt, Archived: ts.Archived,
	}
}

func freeLeave(ts module.TeamScope, now time.Time) bool {
	if ts.FreeLeaveUntil != nil {
		return now.Before(*ts.FreeLeaveUntil)
	}
	if ts.DueAt != nil {
		return now.Before(*ts.DueAt)
	}
	return true
}
