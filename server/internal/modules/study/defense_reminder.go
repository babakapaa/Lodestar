package study

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

const defenseRemindBefore = 2 * time.Hour

type DefenseReminder struct {
	pool     *pgxpool.Pool
	tenants  database.TenantTx
	settings module.Settings
	notifier module.Notifier
	upcoming UpcomingStore
	every    time.Duration
	now      func() time.Time
}

func NewDefenseReminder(k module.Kernel) DefenseReminder {
	return DefenseReminder{
		pool: k.Pool, tenants: k.Tenants, settings: k.Settings, notifier: k.Notifier,
		every: remindEvery, now: time.Now,
	}
}

func (r DefenseReminder) Run(ctx context.Context) {
	ticker := time.NewTicker(r.every)
	defer ticker.Stop()
	for {
		if err := r.RemindAll(ctx); err != nil && ctx.Err() == nil {
			slog.Warn("remind defenses", "error", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (r DefenseReminder) RemindAll(ctx context.Context) error {
	rows, err := r.pool.Query(ctx, `SELECT id::text FROM core.tenants`)
	if err != nil {
		return err
	}
	tenants, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return err
	}
	for _, tenantID := range tenants {
		enabled, err := r.settings.ModuleEnabled(ctx, tenantID, Name)
		if err != nil {
			return err
		}
		if !enabled {
			continue
		}
		if err := r.tenants.Run(ctx, tenantID, func(tx pgx.Tx) error {
			return r.remind(ctx, tx, tenantID)
		}); err != nil {
			return err
		}
	}
	return nil
}

func (r DefenseReminder) remind(ctx context.Context, tx pgx.Tx, tenantID string) error {
	now := r.now()
	list, err := r.upcoming.Between(ctx, tx, now, now.Add(defenseRemindBefore))
	if err != nil {
		return err
	}
	for _, u := range list {
		if err := r.notifier.Notify(ctx, tx, tenantID, module.Notification{
			Type: NotifyDefenseSoon, UserIDs: []string{u.StudentID},
			Params: map[string]any{"assignment": u.Assignment, "course": u.Course},
			Link:   assignmentLink(u.AssignmentID),
		}); err != nil {
			return err
		}
		if err := r.upcoming.Mark(ctx, tx, u); err != nil {
			return err
		}
	}
	return nil
}
