package sandbox

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
)

type RunTable struct {
	Header []string
	Rows   [][]any
}

func (s Service) ExportRuns(ctx context.Context, p auth.Principal) (RunTable, error) {
	out := RunTable{Header: []string{"queued_at", "started_at", "ended_at", "title", "kind", "class", "template", "team", "work", "status", "reason", "hours", "refunded", "gpu_load", "command"}}
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		list, err := s.runs.History(ctx, tx, p.UserID)
		if err != nil {
			return err
		}
		for _, r := range list {
			out.Rows = append(out.Rows, []any{
				stamp(&r.QueuedAt), stamp(r.StartedAt), stamp(r.EndedAt), r.Title, r.Kind, r.ClassName, r.TemplateName, r.TeamName,
				r.WorkTitle, r.Status, r.Reason, number(r.Charged), r.Refunded, number(r.GPULoad), r.Command,
			})
		}
		return nil
	})
	return out, err
}

func stamp(t *time.Time) any {
	if t == nil {
		return nil
	}
	return t.UTC().Format(time.RFC3339)
}

func number(v *float64) any {
	if v == nil {
		return nil
	}
	return *v
}
