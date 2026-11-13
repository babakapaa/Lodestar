package sandbox

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/module"
)

const refundWindow = 5 * time.Minute

func (s Service) end(ctx context.Context, tx pgx.Tx, r Run, e runEnd) error {
	t, err := s.tariff(ctx, tx)
	if err != nil {
		return err
	}
	e.charged = r.cost(e.at, t)
	if e.status == RunFailed && r.StartedAt != nil && r.used(e.at) < refundWindow {
		e.charged, e.refunded = 0, true
	}
	if err := s.runs.End(ctx, tx, r.ID, e); err != nil {
		return err
	}
	if e.charged > 0 {
		if err := s.ledger.Charge(ctx, tx, r, e.charged); err != nil {
			return err
		}
	}
	tenantID, err := tenantOf(ctx, tx)
	if err != nil {
		return err
	}
	if err := s.announce(ctx, tx, tenantID, r, e); err != nil {
		return err
	}
	if e.charged == 0 {
		return nil
	}
	return s.warnLowQuota(ctx, tx, tenantID, r, e.charged)
}

func (s Service) announce(ctx context.Context, tx pgx.Tx, tenantID string, r Run, e runEnd) error {
	reason := e.reason
	if reason == "" {
		reason = r.Reason
	}
	if e.status == RunStopped && (reason == "" || reason == ReasonIdle) {
		return nil
	}
	kind := NotifyRunFailed
	if e.status == RunFinished {
		kind = NotifyRunFinished
	}
	return s.notifier.Notify(ctx, tx, tenantID, module.Notification{
		Type: kind, UserIDs: []string{r.UserID}, Link: runLink(r.ID),
		Params: map[string]any{"title": r.Title, "hours": format(e.charged)},
	})
}

func tenantOf(ctx context.Context, tx pgx.Tx) (string, error) {
	var id string
	err := tx.QueryRow(ctx, `SELECT current_setting('lodestar.tenant')`).Scan(&id)
	return id, err
}

func runLink(id string) string {
	return LinkQuota + "/runs/" + id
}
