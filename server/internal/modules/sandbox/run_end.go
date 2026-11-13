package sandbox

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

func (s Service) end(ctx context.Context, tx pgx.Tx, r Run, status string, exitCode *int, reason string, at time.Time) error {
	charged := r.cost(at)
	if err := s.runs.End(ctx, tx, r.ID, status, exitCode, reason, charged, at); err != nil {
		return err
	}
	if charged == 0 {
		return nil
	}
	return s.ledger.Charge(ctx, tx, r, charged)
}
