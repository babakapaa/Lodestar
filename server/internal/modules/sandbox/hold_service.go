package sandbox

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/audit"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
)

func (s Service) HeldRuns(ctx context.Context, p auth.Principal) ([]Run, error) {
	var out []Run
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if err := s.manage(ctx, tx, p); err != nil {
			return err
		}
		list, err := s.runs.Held(ctx, tx)
		if err != nil {
			return err
		}
		out = make([]Run, 0, len(list))
		for _, r := range list {
			out = append(out, r.seenBy(p.UserID))
		}
		return nil
	})
	return out, err
}

func (s Service) DecideHold(ctx context.Context, p auth.Principal, id string, d HoldDecision) (Run, error) {
	var r Run
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if err := s.manage(ctx, tx, p); err != nil {
			return err
		}
		current, err := s.runs.ByID(ctx, tx, id, true)
		if err != nil {
			return err
		}
		if current.Status != RunHeld {
			return ErrRunNotHeld
		}
		if d.Release {
			err = s.runs.Release(ctx, tx, id)
		} else {
			err = s.runs.Stopping(ctx, tx, id, ReasonWatchdog)
		}
		if err != nil {
			return err
		}
		if err := s.log.Record(ctx, tx, p.UserID, audit.Change{Action: ActionHold, TargetKind: "sandbox.run", TargetID: id, Details: d}); err != nil {
			return err
		}
		r, err = s.view(ctx, tx, id)
		r = r.seenBy(p.UserID)
		return err
	})
	return r, err
}
