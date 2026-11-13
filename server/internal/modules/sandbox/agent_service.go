package sandbox

import (
	"context"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

func (s Service) Sync(ctx context.Context, raw string, report AgentReport) (AgentPlan, error) {
	token, err := ParseNodeToken(raw)
	if err != nil {
		return AgentPlan{}, err
	}
	enabled, err := s.settings.ModuleEnabled(ctx, token.TenantID, Name)
	if err != nil {
		return AgentPlan{}, err
	}
	if !enabled {
		return AgentPlan{}, ErrBadNodeToken
	}
	var plan AgentPlan
	err = s.tenants.Run(ctx, token.TenantID, func(tx pgx.Tx) error {
		n, err := s.nodes.ByToken(ctx, tx, token.Hash())
		if err != nil {
			return err
		}
		if err := s.nodes.Seen(ctx, tx, n.ID); err != nil {
			return err
		}
		runs, err := s.runs.OnNode(ctx, tx, n.ID)
		if err != nil {
			return err
		}
		if err := s.apply(ctx, tx, token.TenantID, runs, report); err != nil {
			return err
		}
		if err := s.Schedule(ctx, tx); err != nil {
			return err
		}
		if runs, err = s.runs.OnNode(ctx, tx, n.ID); err != nil {
			return err
		}
		plan = AgentPlan{NodeID: n.ID, Runs: make([]RunSpec, 0, len(runs))}
		for _, r := range runs {
			plan.Runs = append(plan.Runs, RunSpec{
				RunID: r.ID, Kind: r.Kind, Title: r.Title, Image: r.Image, Command: r.Command,
				MaxHours: r.MaxHours, Workspace: r.workspace(), Hold: r.Status == RunHeld, Stop: r.Status == RunStopping,
			})
		}
		return nil
	})
	return plan, err
}

func (s Service) apply(ctx context.Context, tx pgx.Tx, tenantID string, runs []Run, report AgentReport) error {
	seen := make(map[string]RunReport, len(report.Runs))
	for _, rr := range report.Runs {
		seen[rr.RunID] = rr
	}
	now := time.Now()
	for _, r := range runs {
		rr, ok := seen[r.ID]
		switch {
		case !ok && r.Status == RunStarting:
		case !ok && r.Status == RunStopping:
			if err := s.end(ctx, tx, r, runEnd{status: RunStopped, at: now}); err != nil {
				return err
			}
		case !ok:
			if err := s.end(ctx, tx, r, runEnd{status: RunFailed, reason: ReasonVanished, at: now}); err != nil {
				return err
			}
		case rr.Phase == PhaseRunning:
			if err := s.runs.Started(ctx, tx, r.ID, rr.URL); err != nil {
				return err
			}
			if r.Status == RunStarting {
				r.Status, r.ActiveAt = RunRunning, nil
			}
			if rr.Usage != nil {
				if err := s.observe(ctx, tx, tenantID, r, *rr.Usage, now); err != nil {
					return err
				}
			}
		case rr.Phase == PhaseExited:
			if r.StartedAt == nil {
				started := now
				r.StartedAt = &started
			}
			code := rr.ExitCode
			if err := s.end(ctx, tx, r, runEnd{status: exitStatus(r, code), exitCode: &code, log: tail(rr.Log), at: now}); err != nil {
				return err
			}
		}
	}
	return nil
}

func tail(log string) string {
	const limit = 8000
	runes := []rune(strings.ToValidUTF8(log, ""))
	if len(runes) > limit {
		runes = runes[len(runes)-limit:]
	}
	return string(runes)
}

func exitStatus(r Run, code int) string {
	switch {
	case r.Status == RunStopping:
		return RunStopped
	case code == 0:
		return RunFinished
	default:
		return RunFailed
	}
}
