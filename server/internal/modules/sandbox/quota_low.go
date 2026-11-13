package sandbox

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/module"
)

const lowShare = 0.1

func (s Service) warnLowQuota(ctx context.Context, tx pgx.Tx, tenantID string, r Run, charged float64) error {
	t, err := s.terms.ByID(ctx, tx, r.TermID)
	if err != nil {
		return err
	}
	kind, id := r.owner()
	recipients := []string{r.UserID}
	var lines []QuotaLine
	if kind == OwnerTeam {
		if lines, err = s.ledger.Lines(ctx, tx, t.ID, OwnerTeam, id); err != nil {
			return err
		}
		ref, err := s.team(ctx, tx, r.UserID, id)
		switch {
		case err == nil:
			recipients = ref.Members
		case !errors.Is(err, ErrNotInTeam):
			return err
		}
	} else if lines, err = s.userLines(ctx, tx, t, r.UserID); err != nil {
		return err
	}
	var granted float64
	for _, l := range lines {
		if l.Hours > 0 {
			granted += l.Hours
		}
	}
	left, limit := total(lines), granted*lowShare
	if granted <= 0 || left >= limit || left+charged < limit {
		return nil
	}
	return s.notifier.Notify(ctx, tx, tenantID, module.Notification{
		Type: NotifyQuotaLow, UserIDs: recipients, Link: LinkQuota,
		Params: map[string]any{"hours": format(max(left, 0)), "total": format(hours(granted))},
	})
}
