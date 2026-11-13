package setup

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

type Reminder struct {
	pool     *pgxpool.Pool
	tenants  database.TenantTx
	notifier module.Notifier
	store    Store
	every    time.Duration
	now      func() time.Time
}

func NewReminder(pool *pgxpool.Pool, notifier module.Notifier, every time.Duration) Reminder {
	return Reminder{pool: pool, tenants: database.NewTenantTx(pool), notifier: notifier, every: every, now: time.Now}
}

func (r Reminder) Run(ctx context.Context) {
	ticker := time.NewTicker(r.every)
	defer ticker.Stop()
	for {
		if err := r.RemindAll(ctx); err != nil && ctx.Err() == nil {
			slog.Warn("remind setup due dates", "error", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (r Reminder) RemindAll(ctx context.Context) error {
	rows, err := r.pool.Query(ctx, `SELECT id::text FROM core.tenants`)
	if err != nil {
		return err
	}
	tenants, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return err
	}
	today := r.now().UTC()
	for _, tenantID := range tenants {
		if err := r.tenants.Run(ctx, tenantID, func(tx pgx.Tx) error {
			return r.remind(ctx, tx, tenantID, today)
		}); err != nil {
			return err
		}
	}
	return nil
}

func (r Reminder) remind(ctx context.Context, tx pgx.Tx, tenantID string, today time.Time) error {
	due, err := r.store.DueSoon(ctx, tx, today.Format(time.DateOnly), today.AddDate(0, 0, 1).Format(time.DateOnly))
	if err != nil {
		return err
	}
	for _, sec := range due {
		if err := r.notifier.Notify(ctx, tx, tenantID, module.Notification{
			Type: NotifyDueSoon, UserIDs: []string{*sec.AssigneeID}, Params: map[string]any{"due": *sec.DueOn}, Link: sectionLink(sec.Key),
		}); err != nil {
			return err
		}
		if err := r.store.MarkReminded(ctx, tx, sec.Key, *sec.DueOn); err != nil {
			return err
		}
	}
	return nil
}
