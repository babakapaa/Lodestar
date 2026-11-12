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

const (
	remindEvery  = 15 * time.Minute
	remindBefore = 24 * time.Hour
)

type DeadlineReminder struct {
	pool        *pgxpool.Pool
	tenants     database.TenantTx
	gate        CourseGate
	settings    module.Settings
	notifier    module.Notifier
	assignments AssignmentStore
	courses     CourseStore
	reminders   ReminderStore
	every       time.Duration
	now         func() time.Time
}

func NewDeadlineReminder(k module.Kernel) DeadlineReminder {
	return DeadlineReminder{
		pool: k.Pool, tenants: k.Tenants, gate: CourseGate{access: k.Access}, settings: k.Settings,
		notifier: k.Notifier, every: remindEvery, now: time.Now,
	}
}

func (r DeadlineReminder) Run(ctx context.Context) {
	ticker := time.NewTicker(r.every)
	defer ticker.Stop()
	for {
		if err := r.RemindAll(ctx); err != nil && ctx.Err() == nil {
			slog.Warn("remind deadlines", "error", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (r DeadlineReminder) RemindAll(ctx context.Context) error {
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

func (r DeadlineReminder) remind(ctx context.Context, tx pgx.Tx, tenantID string) error {
	now := r.now()
	due, err := r.assignments.DueBetween(ctx, tx, now, now.Add(remindBefore))
	if err != nil {
		return err
	}
	for _, a := range due {
		students, err := r.gate.StudentIDs(ctx, tx, a.CourseID)
		if err != nil {
			return err
		}
		pending, err := r.reminders.Pending(ctx, tx, a.ID, students)
		if err != nil {
			return err
		}
		if len(pending) == 0 {
			continue
		}
		c, err := r.courses.ByID(ctx, tx, a.CourseID)
		if err != nil {
			return err
		}
		if err := r.notifier.Notify(ctx, tx, tenantID, module.Notification{
			Type: NotifyDeadline, UserIDs: pending,
			Params: map[string]any{"assignment": a.Title, "course": c.Title},
			Link:   assignmentLink(a.ID),
		}); err != nil {
			return err
		}
		if err := r.reminders.Mark(ctx, tx, a.ID, pending); err != nil {
			return err
		}
	}
	return nil
}
