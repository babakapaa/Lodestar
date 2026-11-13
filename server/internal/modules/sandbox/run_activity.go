package sandbox

import (
	"context"
	"math"
	"slices"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/module"
)

const (
	idleWarnAfter = 20 * time.Minute
	idleStopAfter = 30 * time.Minute
)

func (s Service) observe(ctx context.Context, tx pgx.Tx, tenantID string, r Run, u RunUsage, now time.Time) error {
	if r.Status != RunRunning {
		return nil
	}
	watched, err := s.watched(ctx, tx)
	if err != nil {
		return err
	}
	if hit := suspicious(u.Processes, watched); hit != "" {
		return s.runs.Hold(ctx, tx, r.ID, hit)
	}
	active := r.ActiveAt == nil || u.active(r.NetBytes)
	if err := s.runs.Observe(ctx, tx, r.ID, active, u); err != nil {
		return err
	}
	if active || !r.interactive() {
		return nil
	}
	idle := now.Sub(*r.ActiveAt)
	switch {
	case idle >= idleStopAfter:
		return s.runs.Stopping(ctx, tx, r.ID, ReasonIdle)
	case idle >= idleWarnAfter && !r.IdleWarned:
		if err := s.runs.Warned(ctx, tx, r.ID); err != nil {
			return err
		}
		return s.notifier.Notify(ctx, tx, tenantID, module.Notification{
			Type: NotifyIdle, UserIDs: []string{r.UserID}, Link: runLink(r.ID),
			Params: map[string]any{"title": r.Title, "minutes": int(math.Ceil((idleStopAfter - idle).Minutes()))},
		})
	}
	return nil
}

func suspicious(processes, watched []string) string {
	for _, p := range processes {
		name := strings.ToLower(strings.TrimSpace(p))
		if i := strings.LastIndex(name, "/"); i >= 0 {
			name = name[i+1:]
		}
		if name != "" && slices.Contains(watched, name) {
			return name
		}
	}
	return ""
}
