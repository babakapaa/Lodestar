package sandbox

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
)

func (s Service) StartRun(ctx context.Context, p auth.Principal, in RunInput) (Run, error) {
	var r Run
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		var err error
		r, err = s.start(ctx, tx, p, in)
		return err
	})
	return r, err
}

func (s Service) RerunRun(ctx context.Context, p auth.Principal, id string) (Run, error) {
	var r Run
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		old, err := s.runs.ByID(ctx, tx, id, false)
		if err != nil {
			return err
		}
		if old.UserID != p.UserID {
			return ErrRunNotFound
		}
		if old.EndedAt == nil {
			return ErrRunActive
		}
		r, err = s.start(ctx, tx, p, old.input())
		return err
	})
	return r, err
}

func (s Service) start(ctx context.Context, tx pgx.Tx, p auth.Principal, in RunInput) (Run, error) {
	in, err := in.normalized()
	if err != nil {
		return Run{}, err
	}
	t, err := s.current(ctx, tx)
	if err != nil {
		return Run{}, err
	}
	n := newRun{termID: t.ID, userID: p.UserID, in: in}
	if err := s.environment(ctx, tx, &n); err != nil {
		return Run{}, err
	}
	if err := s.workOf(ctx, tx, p, in, &n); err != nil {
		return Run{}, err
	}
	ownerKind, ownerID := OwnerUser, p.UserID
	if in.TeamID != "" {
		ref, err := s.team(ctx, tx, p.UserID, in.TeamID)
		if err != nil {
			return Run{}, err
		}
		n.teamID, n.teamName = &ref.ID, ref.Name
		ownerKind, ownerID = OwnerTeam, ref.ID
	}
	if n.weight, err = s.weight(ctx, tx, t, p.UserID); err != nil {
		return Run{}, err
	}
	if err := s.afford(ctx, tx, t, ownerKind, ownerID, p.UserID, n.in.cost(n.rate)); err != nil {
		return Run{}, err
	}
	id, err := s.runs.Insert(ctx, tx, n)
	if err != nil {
		return Run{}, err
	}
	return s.view(ctx, tx, id)
}

func (s Service) environment(ctx context.Context, tx pgx.Tx, n *newRun) error {
	in := n.in
	tpl, err := s.templates.ByID(ctx, tx, in.TemplateID)
	if err != nil {
		return err
	}
	if tpl.Archived {
		return ErrTemplateNotFound
	}
	classes, err := s.classes.All(ctx, tx)
	if err != nil {
		return err
	}
	i := slices.IndexFunc(classes, func(c GPUClass) bool { return c.ID == in.ClassID && !c.Archived })
	if i < 0 {
		return ErrClassNotFound
	}
	n.rate = classes[i].Rate
	if !tpl.Starter {
		return nil
	}
	if in.TeamID != "" {
		return ErrInvalidRun
	}
	busy, err := s.runs.StarterActive(ctx, tx, n.userID)
	if err != nil {
		return err
	}
	if busy {
		return ErrStarterBusy
	}
	n.rate, n.in.MaxHours = 0, min(in.MaxHours, starterHours)
	return nil
}

const starterHours = 0.25

func (s Service) workOf(ctx context.Context, tx pgx.Tx, p auth.Principal, in RunInput, n *newRun) error {
	if in.WorkKind == "" {
		return nil
	}
	b, err := s.brief(ctx, tx, p, in.WorkKind, in.WorkID)
	if err != nil {
		return err
	}
	n.workTitle, n.dueAt = b.Title, b.DueAt
	return nil
}

func (s Service) weight(ctx context.Context, tx pgx.Tx, t Term, userID string) (float64, error) {
	from, to := t.bounds()
	grants, err := s.grants.Grants(ctx, tx, userID, from, to)
	if err != nil {
		return 0, err
	}
	var w float64
	for _, g := range grants {
		if g.Priority {
			w = max(w, g.Hours)
		}
	}
	return w, nil
}

func (s Service) afford(ctx context.Context, tx pgx.Tx, t Term, ownerKind, ownerID, userID string, cost float64) error {
	if err := s.ledger.Lock(ctx, tx, t.ID, ownerKind, ownerID); err != nil {
		return err
	}
	var lines []QuotaLine
	var err error
	if ownerKind == OwnerTeam {
		lines, err = s.ledger.Lines(ctx, tx, t.ID, OwnerTeam, ownerID)
	} else {
		lines, err = s.userLines(ctx, tx, t, userID)
	}
	if err != nil {
		return err
	}
	reserved, err := s.runs.Reserved(ctx, tx, t.ID, ownerKind, ownerID)
	if err != nil {
		return err
	}
	if total(lines)-reserved < cost {
		return ErrNotEnoughHours
	}
	return nil
}

func (s Service) Runs(ctx context.Context, p auth.Principal) ([]Run, error) {
	out := []Run{}
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		t, err := s.current(ctx, tx)
		if errors.Is(err, ErrNoTerm) {
			return nil
		}
		if err != nil {
			return err
		}
		teamIDs, err := s.teamIDs(ctx, tx, p.UserID)
		if err != nil {
			return err
		}
		list, err := s.runs.Visible(ctx, tx, t.ID, p.UserID, teamIDs)
		if err != nil {
			return err
		}
		if out, err = s.placed(ctx, tx, list); err != nil {
			return err
		}
		for i := range out {
			out[i] = out[i].seenBy(p.UserID)
			out[i].Log = ""
		}
		return nil
	})
	return out, err
}

func (s Service) GetRun(ctx context.Context, p auth.Principal, id string) (Run, error) {
	var r Run
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		var err error
		if r, err = s.visible(ctx, tx, p, id, false); err != nil {
			return err
		}
		r, err = s.view(ctx, tx, id)
		r = r.seenBy(p.UserID)
		return err
	})
	return r, err
}

func (s Service) StopRun(ctx context.Context, p auth.Principal, id string) (Run, error) {
	var r Run
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		current, err := s.visible(ctx, tx, p, id, true)
		if err != nil {
			return err
		}
		if current.UserID != p.UserID {
			return ErrRunNotFound
		}
		switch current.Status {
		case RunQueued:
			err = s.runs.End(ctx, tx, id, runEnd{status: RunStopped, at: time.Now()})
		case RunStarting, RunRunning:
			err = s.runs.Stopping(ctx, tx, id, "")
		case RunHeld:
			err = s.runs.Stopping(ctx, tx, id, ReasonWatchdog)
		default:
			return ErrRunEnded
		}
		if err != nil {
			return err
		}
		r, err = s.view(ctx, tx, id)
		r = r.seenBy(p.UserID)
		return err
	})
	return r, err
}

func (s Service) visible(ctx context.Context, tx pgx.Tx, p auth.Principal, id string, lock bool) (Run, error) {
	r, err := s.runs.ByID(ctx, tx, id, lock)
	if err != nil {
		return r, err
	}
	if r.UserID == p.UserID {
		return r, nil
	}
	if r.TeamID != nil {
		teamIDs, err := s.teamIDs(ctx, tx, p.UserID)
		if err != nil {
			return r, err
		}
		if slices.Contains(teamIDs, *r.TeamID) {
			return r, nil
		}
	}
	return Run{}, ErrRunNotFound
}

func (s Service) view(ctx context.Context, tx pgx.Tx, id string) (Run, error) {
	r, err := s.runs.ByID(ctx, tx, id, false)
	if err != nil {
		return r, err
	}
	list, err := s.placed(ctx, tx, []Run{r})
	if err != nil {
		return r, err
	}
	return list[0], nil
}

func (s Service) placed(ctx context.Context, tx pgx.Tx, list []Run) ([]Run, error) {
	if !slices.ContainsFunc(list, func(r Run) bool { return r.Status == RunQueued }) {
		return list, nil
	}
	days, err := s.integer(ctx, tx, DeadlineDays)
	if err != nil {
		return nil, err
	}
	positions, err := s.runs.Positions(ctx, tx, days)
	if err != nil {
		return nil, err
	}
	for i, r := range list {
		if pos, ok := positions[r.ID]; ok {
			list[i].Position = &pos
		}
	}
	return list, nil
}

func (s Service) teamIDs(ctx context.Context, tx pgx.Tx, userID string) ([]string, error) {
	if s.teams == nil {
		return []string{}, nil
	}
	refs, err := s.teams.Teams(ctx, tx, userID)
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(refs))
	for _, r := range refs {
		ids = append(ids, r.ID)
	}
	return ids, nil
}
