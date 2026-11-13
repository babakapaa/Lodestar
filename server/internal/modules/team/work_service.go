package team

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/module"
)

func (s Service) Propose(ctx context.Context, p auth.Principal, id string, in ProposalInput) (TeamDetails, error) {
	in, err := in.normalized()
	if err != nil {
		return TeamDetails{}, err
	}
	return s.change(ctx, p, id, func(tx pgx.Tx, t Team, ts module.TeamScope, _ ScopeRole) error {
		if !t.has(p.UserID) {
			return ErrTeamNotFound
		}
		if ts.Archived || t.ArchivedAt != nil {
			return ErrScopeArchived
		}
		source, ok := s.scopes.Find(t.ScopeKind)
		if !ok {
			return ErrScopeNotFound
		}
		work, err := source.PrepareTeamWork(ctx, tx, t.ScopeID, p, in.Work)
		if err != nil {
			return err
		}
		now := time.Now()
		proposalID, err := s.proposals.Insert(ctx, tx, t.ID, p.UserID, work, in.Note, silentAt(now, ts.DueAt), ts.DueAt)
		if err != nil {
			return err
		}
		if err := s.proposals.Vote(ctx, tx, proposalID, p.UserID, true, ""); err != nil {
			return err
		}
		if err := s.settle(ctx, tx, p.TenantID, proposalID, now); err != nil {
			return err
		}
		others := otherThan(t.ids(), p.UserID)
		if len(others) == 0 {
			return nil
		}
		author, err := s.name(ctx, tx, p.UserID)
		if err != nil {
			return err
		}
		return s.notifier.Notify(ctx, tx, p.TenantID, module.Notification{
			Type: NotifyProposed, UserIDs: others,
			Params: map[string]any{"scope": ts.Title, "team": t.Name, "author": author},
			Link:   teamLink(t.ID),
		})
	})
}

func (s Service) Vote(ctx context.Context, p auth.Principal, proposalID string, in VoteInput) (TeamDetails, error) {
	in, err := in.normalized()
	if err != nil {
		return TeamDetails{}, err
	}
	var teamID string
	err = s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		proposal, err := s.proposals.ByID(ctx, tx, proposalID, false)
		teamID = proposal.TeamID
		return err
	})
	if err != nil {
		return TeamDetails{}, err
	}
	return s.change(ctx, p, teamID, func(tx pgx.Tx, t Team, ts module.TeamScope, _ ScopeRole) error {
		if !t.has(p.UserID) {
			return ErrProposalNotFound
		}
		proposal, err := s.proposals.ByID(ctx, tx, proposalID, true)
		if err != nil {
			return err
		}
		if proposal.Status != ProposalOpen {
			return ErrProposalNotFound
		}
		if err := s.proposals.Vote(ctx, tx, proposalID, p.UserID, in.Agree, in.Comment); err != nil {
			return err
		}
		if !in.Agree {
			author, err := s.name(ctx, tx, p.UserID)
			if err != nil {
				return err
			}
			if err := s.notifier.Notify(ctx, tx, p.TenantID, module.Notification{
				Type: NotifyObjected, UserIDs: otherThan(t.ids(), p.UserID),
				Params: map[string]any{"scope": ts.Title, "team": t.Name, "author": author},
				Link:   teamLink(t.ID),
			}); err != nil {
				return err
			}
		}
		return s.settle(ctx, tx, p.TenantID, proposalID, time.Now())
	})
}

func (s Service) settle(ctx context.Context, tx pgx.Tx, tenantID, proposalID string, now time.Time) error {
	proposal, err := s.proposals.ByID(ctx, tx, proposalID, false)
	if err != nil {
		return err
	}
	t, err := s.teams.ByID(ctx, tx, proposal.TeamID, true)
	if err != nil {
		return err
	}
	if proposal, err = s.proposals.ByID(ctx, tx, proposalID, true); err != nil {
		return err
	}
	if proposal.Votes, err = s.proposals.Votes(ctx, tx, proposal.ID, t.Members); err != nil {
		return err
	}
	if !proposal.ready(now) {
		return nil
	}
	source, ok := s.scopes.Find(t.ScopeKind)
	if !ok {
		return ErrScopeNotFound
	}
	if err := source.DeliverTeamWork(ctx, tx, proposal.delivery(tenantID, t)); err != nil {
		return err
	}
	if err := s.proposals.Sent(ctx, tx, proposal.ID); err != nil {
		return err
	}
	ts, err := s.scopes.Load(ctx, tx, t.ScopeKind, t.ScopeID)
	if err != nil {
		return err
	}
	return s.notifier.Notify(ctx, tx, tenantID, module.Notification{
		Type: NotifySubmitted, UserIDs: t.ids(),
		Params: map[string]any{"scope": ts.Title, "team": t.Name},
		Link:   teamLink(t.ID),
	})
}

func (s Service) SettleDue(ctx context.Context, tenantID string, now time.Time) error {
	var due []string
	err := s.tenants.Run(ctx, tenantID, func(tx pgx.Tx) error {
		var err error
		due, err = s.proposals.Due(ctx, tx, now)
		return err
	})
	if err != nil {
		return err
	}
	for _, id := range due {
		err := s.tenants.Run(ctx, tenantID, func(tx pgx.Tx) error {
			return s.settle(ctx, tx, tenantID, id, now)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func otherThan(ids []string, userID string) []string {
	out := make([]string, 0, len(ids))
	for _, id := range ids {
		if id != userID {
			out = append(out, id)
		}
	}
	return out
}
