package sandbox

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

func (s Service) Schedule(ctx context.Context, tx pgx.Tx) error {
	if _, err := tx.Exec(ctx, `SELECT pg_advisory_xact_lock(hashtext('sandbox.schedule:' || current_setting('lodestar.tenant')))`); err != nil {
		return err
	}
	if err := s.reapLost(ctx, tx); err != nil {
		return err
	}
	overdue, err := s.runs.Overdue(ctx, tx)
	if err != nil {
		return err
	}
	for _, r := range overdue {
		if err := s.runs.Stopping(ctx, tx, r.ID, ReasonTimeLimit); err != nil {
			return err
		}
	}
	days, err := s.integer(ctx, tx, DeadlineDays)
	if err != nil {
		return err
	}
	classes, err := s.runs.QueuedClasses(ctx, tx)
	if err != nil {
		return err
	}
	for _, classID := range classes {
		if err := s.fill(ctx, tx, classID, days); err != nil {
			return err
		}
	}
	return nil
}

func (s Service) fill(ctx context.Context, tx pgx.Tx, classID string, days int) error {
	nodes, err := s.nodes.Free(ctx, tx, classID)
	if err != nil || len(nodes) == 0 {
		return err
	}
	slots := 0
	for _, n := range nodes {
		slots += n.GPUs - n.Busy
	}
	queued, err := s.runs.Queued(ctx, tx, classID, days, slots)
	if err != nil {
		return err
	}
	i := 0
	for _, r := range queued {
		for nodes[i].Busy >= nodes[i].GPUs {
			i++
		}
		if err := s.runs.Assign(ctx, tx, r.ID, nodes[i].ID); err != nil {
			return err
		}
		nodes[i].Busy++
	}
	return nil
}

func (s Service) reapLost(ctx context.Context, tx pgx.Tx) error {
	lost, err := s.nodes.Lost(ctx, tx)
	if err != nil {
		return err
	}
	for _, n := range lost {
		runs, err := s.runs.OnNode(ctx, tx, n.ID)
		if err != nil {
			return err
		}
		at := time.Now()
		if n.SeenAt != nil {
			at = *n.SeenAt
		}
		for _, r := range runs {
			if err := s.end(ctx, tx, r, RunFailed, nil, ReasonNodeLost, at); err != nil {
				return err
			}
		}
	}
	return nil
}
