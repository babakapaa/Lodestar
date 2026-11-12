package privacy

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/babakapaa/lodestar/server/internal/core/audit"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

const ActionAccountDeleted = "core.account.deleted"

type AccountEraser struct {
	pool      *pgxpool.Pool
	tenants   database.TenantTx
	sections  Sections
	every     time.Duration
	now       func() time.Time
	deletions DeletionStore
	log       audit.Log
}

func NewAccountEraser(pool *pgxpool.Pool, sections Sections, every time.Duration) AccountEraser {
	return AccountEraser{pool: pool, tenants: database.NewTenantTx(pool), sections: sections, every: every, now: time.Now}
}

func (e AccountEraser) WithClock(now func() time.Time) AccountEraser {
	e.now = now
	return e
}

func (e AccountEraser) Run(ctx context.Context) {
	ticker := time.NewTicker(e.every)
	defer ticker.Stop()
	for {
		if err := e.EraseDue(ctx); err != nil && ctx.Err() == nil {
			slog.Warn("erase accounts", "error", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (e AccountEraser) EraseDue(ctx context.Context) error {
	rows, err := e.pool.Query(ctx, `SELECT id::text FROM core.tenants`)
	if err != nil {
		return err
	}
	ids, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return err
	}
	for _, tenantID := range ids {
		var due []string
		err := e.tenants.Run(ctx, tenantID, func(tx pgx.Tx) error {
			var err error
			due, err = e.deletions.Due(ctx, tx, e.now())
			return err
		})
		if err != nil {
			return err
		}
		for _, userID := range due {
			if err := e.erase(ctx, tenantID, userID); err != nil {
				return err
			}
		}
	}
	return nil
}

func (e AccountEraser) erase(ctx context.Context, tenantID, userID string) error {
	return e.tenants.Run(ctx, tenantID, func(tx pgx.Tx) error {
		found, err := e.deletions.Delete(ctx, tx, userID)
		if err != nil || !found {
			return err
		}
		for _, s := range e.sections {
			if eraser, ok := s.(Eraser); ok {
				if err := eraser.Erase(ctx, tx, tenantID, userID); err != nil {
					return fmt.Errorf("erase %s: %w", s.Name(), err)
				}
			}
		}
		if _, err := tx.Exec(ctx, `DELETE FROM core.users WHERE id::text = $1`, userID); err != nil {
			return fmt.Errorf("delete account: %w", err)
		}
		return e.log.Record(ctx, tx, "", audit.Change{Action: ActionAccountDeleted, TargetKind: "user", TargetID: userID})
	})
}
