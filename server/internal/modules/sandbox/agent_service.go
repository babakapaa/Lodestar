package sandbox

import (
	"context"
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
		if err := s.apply(ctx, tx, runs, report); err != nil {
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
				MaxHours: r.MaxHours, Workspace: r.workspace(), Stop: r.Status == RunStopping,
			})
		}
		return nil
	})
	return plan, err
}

func (s Service) apply(ctx context.Context, tx pgx.Tx, runs []Run, report AgentReport) error {
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
			if err := s.end(ctx, tx, r, RunStopped, nil, "", now); err != nil {
				return err
			}
		case !ok:
			if err := s.end(ctx, tx, r, RunFailed, nil, ReasonVanished, now); err != nil {
				return err
			}
		case rr.Phase == PhaseRunning:
			if err := s.runs.Started(ctx, tx, r.ID, rr.URL); err != nil {
				return err
			}
		case rr.Phase == PhaseExited:
			if r.StartedAt == nil {
				started := now
				r.StartedAt = &started
			}
			code := rr.ExitCode
			if err := s.end(ctx, tx, r, exitStatus(r, code), &code, "", now); err != nil {
				return err
			}
		}
	}
	return nil
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
