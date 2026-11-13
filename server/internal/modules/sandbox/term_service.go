package sandbox

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/audit"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
)

func (s Service) Terms(ctx context.Context, p auth.Principal) ([]Term, error) {
	var list []Term
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		var err error
		list, err = s.terms.All(ctx, tx)
		return err
	})
	return list, err
}

func (s Service) CreateTerm(ctx context.Context, p auth.Principal, in TermInput) (Term, error) {
	return s.saveTerm(ctx, p, "", in)
}

func (s Service) UpdateTerm(ctx context.Context, p auth.Principal, id string, in TermInput) (Term, error) {
	return s.saveTerm(ctx, p, id, in)
}

func (s Service) saveTerm(ctx context.Context, p auth.Principal, id string, in TermInput) (Term, error) {
	in, err := in.normalized()
	if err != nil {
		return Term{}, err
	}
	var t Term
	err = s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if err := s.manage(ctx, tx, p); err != nil {
			return err
		}
		if err := s.terms.Lock(ctx, tx); err != nil {
			return err
		}
		overlaps, err := s.terms.Overlaps(ctx, tx, id, in)
		if err != nil {
			return err
		}
		if overlaps {
			return ErrTermOverlap
		}
		if id == "" {
			t, err = s.terms.Insert(ctx, tx, in)
		} else {
			t, err = s.terms.Update(ctx, tx, id, in)
		}
		if err != nil {
			return err
		}
		return s.log.Record(ctx, tx, p.UserID, audit.Change{Action: ActionTerm, TargetKind: "sandbox.term", TargetID: t.ID, Details: in})
	})
	return t, err
}

func (s Service) DeleteTerm(ctx context.Context, p auth.Principal, id string) error {
	return s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if err := s.manage(ctx, tx, p); err != nil {
			return err
		}
		if err := s.terms.Delete(ctx, tx, id); err != nil {
			return err
		}
		return s.log.Record(ctx, tx, p.UserID, audit.Change{Action: ActionTermGone, TargetKind: "sandbox.term", TargetID: id})
	})
}
