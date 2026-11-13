package sandbox

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/audit"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
)

func (s Service) Transfer(ctx context.Context, p auth.Principal, in TransferInput) (Quota, error) {
	in, err := in.normalized()
	if err != nil {
		return Quota{}, err
	}
	var q Quota
	err = s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		t, err := s.current(ctx, tx)
		if err != nil {
			return err
		}
		ref, err := s.team(ctx, tx, p.UserID, in.TeamID)
		if err != nil {
			return err
		}
		if err := s.ledger.Lock(ctx, tx, t.ID, OwnerUser, p.UserID); err != nil {
			return err
		}
		lines, err := s.userLines(ctx, tx, t, p.UserID)
		if err != nil {
			return err
		}
		if total(lines) < in.Hours {
			return ErrNotEnoughHours
		}
		if err := s.ledger.Add(ctx, tx, t.ID, OwnerUser, p.UserID, EntryTransfer, -in.Hours, ref.Name, p.UserID); err != nil {
			return err
		}
		name, err := s.name(ctx, tx, p.UserID)
		if err != nil {
			return err
		}
		if err := s.ledger.Add(ctx, tx, t.ID, OwnerTeam, ref.ID, EntryTransfer, in.Hours, name, p.UserID); err != nil {
			return err
		}
		if err := s.log.Record(ctx, tx, p.UserID, audit.Change{Action: ActionTransfer, TargetKind: "team", TargetID: ref.ID, Details: in}); err != nil {
			return err
		}
		q, err = s.quota(ctx, tx, t, p.UserID)
		return err
	})
	return q, err
}
