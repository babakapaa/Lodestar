package setup

import (
	"context"
	"slices"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/audit"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/module"
)

const maxPilotCourses = 200

func (s Service) PilotCourses(ctx context.Context, p auth.Principal) ([]module.PilotCourse, error) {
	out := []module.PilotCourse{}
	err := s.run(ctx, p, func(tx pgx.Tx, st State) error {
		if err := s.require(ctx, tx, p, Manage); err != nil {
			return err
		}
		var err error
		out, err = s.pilotCourses(ctx, tx, p.TenantID)
		return err
	})
	return out, err
}

func (s Service) pilotCourses(ctx context.Context, tx pgx.Tx, tenantID string) ([]module.PilotCourse, error) {
	out := []module.PilotCourse{}
	for _, src := range s.pilots {
		list, err := src.PilotCourses(ctx, tx, tenantID)
		if err != nil {
			return nil, err
		}
		out = append(out, list...)
	}
	return out, nil
}

func (s Service) validPilot(ctx context.Context, tx pgx.Tx, tenantID string, ids []string) error {
	if len(ids) > maxPilotCourses {
		return ErrInvalidPilot
	}
	known, err := s.pilotCourses(ctx, tx, tenantID)
	if err != nil {
		return err
	}
	for i, id := range ids {
		found := slices.ContainsFunc(known, func(c module.PilotCourse) bool { return c.ID == id })
		if !found || slices.Contains(ids[:i], id) {
			return ErrInvalidPilot
		}
	}
	return nil
}

func (s Service) inPilot(ctx context.Context, tx pgx.Tx, p auth.Principal, st State) (bool, error) {
	if manages, err := s.allowed(ctx, tx, p, Manage); err != nil || manages {
		return manages, err
	}
	for _, src := range s.pilots {
		in, err := src.InPilot(ctx, tx, p.UserID, st.Pilot)
		if err != nil || in {
			return in, err
		}
	}
	return false, nil
}

func (s Service) Expand(ctx context.Context, p auth.Principal) (Overview, error) {
	err := s.run(ctx, p, func(tx pgx.Tx, st State) error {
		if err := s.require(ctx, tx, p, Manage); err != nil {
			return err
		}
		if st.LaunchedAt == nil || len(st.Pilot) == 0 {
			return ErrNotPilot
		}
		if err := s.store.Expand(ctx, tx); err != nil {
			return err
		}
		return s.log.Record(ctx, tx, p.UserID, audit.Change{Action: ActionExpanded, TargetKind: "setup", TargetID: "launch"})
	})
	if err != nil {
		return Overview{}, err
	}
	return s.Overview(ctx, p)
}
