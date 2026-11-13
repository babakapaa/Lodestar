package sandbox

import (
	"context"
	"errors"
	"strconv"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/audit"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/module"
)

const (
	LinkQuota  = "/sandbox"
	LinkTopups = "/sandbox/topups"
)

func (s Service) RequestTopup(ctx context.Context, p auth.Principal, in TopupInput) (Topup, error) {
	in, err := in.normalized()
	if err != nil {
		return Topup{}, err
	}
	var t Topup
	err = s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		term, err := s.current(ctx, tx)
		if err != nil {
			return err
		}
		b, err := s.brief(ctx, tx, p, in.WorkKind, in.WorkID)
		if err != nil {
			return err
		}
		if b.Teaches {
			return access.ErrForbidden
		}
		if t, err = s.topups.Insert(ctx, tx, term.ID, p.UserID, b.Title, in); err != nil {
			return err
		}
		t.Mine = true
		return s.notifier.Notify(ctx, tx, p.TenantID, module.Notification{
			Type: NotifyTopupRequested, UserIDs: b.Teachers, Link: LinkTopups,
			Params: map[string]any{"student": t.name(), "hours": format(t.Hours), "work": b.Title},
		})
	})
	return t, err
}

func (s Service) Topups(ctx context.Context, p auth.Principal) ([]Topup, error) {
	out := []Topup{}
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		term, err := s.current(ctx, tx)
		if errors.Is(err, ErrNoTerm) {
			return nil
		}
		if err != nil {
			return err
		}
		all, err := s.topups.InTerm(ctx, tx, term.ID)
		if err != nil {
			return err
		}
		teaches := map[string]bool{}
		for _, t := range all {
			if t.UserID == p.UserID {
				t.Mine = true
				out = append(out, t)
				continue
			}
			key := t.WorkKind + "/" + t.WorkID
			ok, seen := teaches[key]
			if !seen {
				b, err := s.briefs.Brief(ctx, tx, p, t.WorkKind, t.WorkID)
				if err != nil && !errors.Is(err, module.ErrBriefNotFound) {
					return err
				}
				ok = err == nil && b.Teaches
				teaches[key] = ok
			}
			if ok {
				out = append(out, t)
			}
		}
		return nil
	})
	return out, err
}

func (s Service) DecideTopup(ctx context.Context, p auth.Principal, id string, d TopupDecision) (Topup, error) {
	var t Topup
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		pending, err := s.topups.Pending(ctx, tx, id)
		if err != nil {
			return err
		}
		b, err := s.brief(ctx, tx, p, pending.WorkKind, pending.WorkID)
		if err != nil {
			return err
		}
		if !b.Teaches {
			return access.ErrForbidden
		}
		granted, err := d.granted(pending)
		if err != nil {
			return err
		}
		status, kind, amount := TopupDeclined, NotifyTopupDeclined, (*float64)(nil)
		if d.Approve {
			status, kind, amount = TopupApproved, NotifyTopupApproved, &granted
			if err := s.ledger.Add(ctx, tx, pending.TermID, OwnerUser, pending.UserID, EntryTopup, granted, pending.Title, p.UserID); err != nil {
				return err
			}
		}
		if t, err = s.topups.Decide(ctx, tx, id, status, amount, p.UserID); err != nil {
			return err
		}
		if err := s.log.Record(ctx, tx, p.UserID, audit.Change{Action: ActionTopup, TargetKind: "sandbox.topup", TargetID: id, Details: d}); err != nil {
			return err
		}
		return s.notifier.Notify(ctx, tx, p.TenantID, module.Notification{
			Type: kind, UserIDs: []string{t.UserID}, Link: LinkQuota,
			Params: map[string]any{"hours": format(granted), "work": t.Title},
		})
	})
	return t, err
}

func (s Service) brief(ctx context.Context, tx pgx.Tx, p auth.Principal, kind, id string) (module.Brief, error) {
	b, err := s.briefs.Brief(ctx, tx, p, kind, id)
	if errors.Is(err, module.ErrBriefNotFound) {
		return b, ErrWorkNotFound
	}
	return b, err
}

func format(h float64) string {
	return strconv.FormatFloat(h, 'f', -1, 64)
}
