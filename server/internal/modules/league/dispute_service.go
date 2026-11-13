package league

import (
	"context"
	"slices"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
)

func (s Service) Dispute(ctx context.Context, p auth.Principal, leagueID string, in DisputeInput) (Dispute, error) {
	in, err := in.normalized()
	if err != nil {
		return Dispute{}, err
	}
	var d Dispute
	err = s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		r, _, err := s.open(ctx, tx, p, leagueID, RoleStudent)
		if err != nil {
			return err
		}
		if err := s.target(ctx, tx, p, leagueID, in); err != nil {
			return err
		}
		if d, err = s.disputes.Insert(ctx, tx, leagueID, p.UserID, in); err != nil {
			return err
		}
		people, err := s.people.ByIDs(ctx, tx, []string{p.UserID})
		if err != nil {
			return err
		}
		author := ""
		if len(people) > 0 {
			author = people[0].DisplayName
		}
		return s.tellDisputed(ctx, tx, p.TenantID, r, author)
	})
	return d, err
}

func (s Service) target(ctx context.Context, tx pgx.Tx, p auth.Principal, leagueID string, in DisputeInput) error {
	if in.EntryID != nil {
		e, err := s.entries.ByID(ctx, tx, leagueID, *in.EntryID)
		if err != nil {
			return err
		}
		if e.UserID != p.UserID {
			return ErrDisputeTarget
		}
		return nil
	}
	t, owner, err := s.transitions.ByID(ctx, tx, *in.TransitionID)
	if err != nil || owner != leagueID || t.UserID != p.UserID || t.Status != TransitionApplied {
		return ErrDisputeTarget
	}
	return nil
}

func (s Service) Answer(ctx context.Context, p auth.Principal, id string, in AnswerInput) (Dispute, error) {
	in, err := in.normalized()
	if err != nil {
		return Dispute{}, err
	}
	var out Dispute
	err = s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		d, leagueID, err := s.disputes.ByID(ctx, tx, id)
		if err != nil {
			return err
		}
		r, source, err := s.open(ctx, tx, p, leagueID, RoleLead)
		if err != nil {
			return ErrDisputeNotFound
		}
		if d.Status != DisputeOpen {
			return ErrDisputeNotFound
		}
		if in.Points != 0 && (d.EntryID == nil || r.Archived) {
			return ErrInvalidInput
		}
		at, err := s.sync(ctx, tx, p.TenantID, r, source)
		if err != nil {
			return err
		}
		status := DisputeRejected
		if in.Accept {
			status = DisputeAccepted
			if err := s.settle(ctx, tx, at, d, in, p.UserID); err != nil {
				return err
			}
		}
		if out, err = s.disputes.Decide(ctx, tx, d.ID, status, in.Answer, p.UserID); err != nil {
			return err
		}
		if err := s.record(ctx, tx, p, ActionDispute, leagueID, map[string]any{"disputeId": d.ID, "accept": in.Accept, "points": in.Points}); err != nil {
			return err
		}
		if in.Points != 0 {
			if err := s.promote(ctx, tx, p.TenantID, at); err != nil {
				return err
			}
		}
		return s.tellAnswered(ctx, tx, p.TenantID, r, out)
	})
	return out, err
}

func (s Service) settle(ctx context.Context, tx pgx.Tx, at synced, d Dispute, in AnswerInput, deciderID string) error {
	if d.EntryID != nil {
		if in.Points == 0 {
			return nil
		}
		_, err := s.entries.Adjust(ctx, tx, at.league.ID, at.season.ID, d.UserID, in.Points, in.Answer, deciderID)
		return err
	}
	t, _, err := s.transitions.ByID(ctx, tx, *d.TransitionID)
	if err != nil {
		return err
	}
	if t.Status != TransitionApplied {
		return nil
	}
	if _, err := s.transitions.Decide(ctx, tx, t.ID, TransitionReverted, deciderID); err != nil {
		return err
	}
	if !slices.Contains(at.roster.Members, t.UserID) {
		return nil
	}
	return s.players.SetDivision(ctx, tx, at.league.ID, t.UserID, t.From)
}
