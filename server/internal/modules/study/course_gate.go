package study

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/structure"
)

type CourseGate struct {
	access    access.Service
	courses   CourseStore
	roster    RosterStore
	directory structure.Directory
}

func (g CourseGate) Visible(ctx context.Context, tx pgx.Tx, p auth.Principal, id string) (Course, Standing, error) {
	c, err := g.courses.ByID(ctx, tx, id)
	if err != nil {
		return Course{}, Standing{}, err
	}
	st, err := g.standing(ctx, tx, p, c)
	if err != nil {
		return Course{}, Standing{}, err
	}
	if !st.Visible() {
		return Course{}, Standing{}, ErrCourseNotFound
	}
	return c, st, nil
}

func (g CourseGate) Managed(ctx context.Context, tx pgx.Tx, p auth.Principal, id string) (Course, error) {
	c, st, err := g.Visible(ctx, tx, p, id)
	if err != nil {
		return Course{}, err
	}
	if !st.Manage {
		return Course{}, access.ErrForbidden
	}
	return c, nil
}

func (g CourseGate) standing(ctx context.Context, tx pgx.Tx, p auth.Principal, c Course) (Standing, error) {
	var st Standing
	err := g.access.Require(ctx, tx, p, g.Scope(p, c.UnitID), ManageCourses)
	switch {
	case err == nil:
		st.Manage = true
	case !isForbidden(err):
		return Standing{}, err
	}
	teaches, err := g.roster.Teaches(ctx, tx, c.ID, p.UserID)
	if err != nil {
		return Standing{}, err
	}
	if teaches {
		st.Role = RoleTeacher
		return st, nil
	}
	units, err := g.directory.UnitsAbove(ctx, tx, p.UserID)
	if err != nil {
		return Standing{}, err
	}
	learns, err := g.roster.Learns(ctx, tx, c.ID, p.UserID, units)
	if err != nil {
		return Standing{}, err
	}
	if learns {
		st.Role = RoleStudent
	}
	return st, nil
}

func (g CourseGate) Scope(p auth.Principal, unitID *string) access.Scope {
	if unitID == nil {
		return access.TenantScope(p.TenantID)
	}
	return structure.UnitScope(*unitID)
}

func (g CourseGate) StudentIDs(ctx context.Context, tx pgx.Tx, courseID string) ([]string, error) {
	units, err := g.roster.Units(ctx, tx, courseID)
	if err != nil {
		return nil, err
	}
	fromUnits, err := g.directory.UsersUnder(ctx, tx, units)
	if err != nil {
		return nil, err
	}
	direct, err := g.roster.Students(ctx, tx, courseID)
	if err != nil {
		return nil, err
	}
	teachers, err := g.roster.Teachers(ctx, tx, courseID)
	if err != nil {
		return nil, err
	}
	skip := make(map[string]bool, len(teachers))
	for _, id := range teachers {
		skip[id] = true
	}
	var all []string
	for _, id := range append(direct, fromUnits...) {
		if !skip[id] {
			skip[id] = true
			all = append(all, id)
		}
	}
	return all, nil
}
