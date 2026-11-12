package study

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/audit"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/structure"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

type Service struct {
	tenants   database.TenantTx
	access    access.Service
	gate      CourseGate
	courses   CourseStore
	roster    RosterStore
	directory structure.Directory
	units     structure.UnitStore
	people    identity.People
	log       audit.Log
}

func NewService(tenants database.TenantTx, acc access.Service) Service {
	return Service{tenants: tenants, access: acc, gate: CourseGate{access: acc}}
}

func (s Service) Mine(ctx context.Context, p auth.Principal, archived bool) ([]MyCourse, error) {
	var list []MyCourse
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		units, err := s.directory.UnitsAbove(ctx, tx, p.UserID)
		if err != nil {
			return err
		}
		list, err = s.courses.Mine(ctx, tx, p.UserID, units, archived)
		return err
	})
	return list, err
}

func (s Service) Catalog(ctx context.Context, p auth.Principal, unitID *string, archived bool) ([]Course, error) {
	var list []Course
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if err := s.access.Require(ctx, tx, p, s.gate.Scope(p, unitID), ManageCourses); err != nil {
			return err
		}
		var err error
		list, err = s.courses.ByUnit(ctx, tx, unitID, archived)
		return err
	})
	return list, err
}

func (s Service) Create(ctx context.Context, p auth.Principal, in CourseInput) (Course, error) {
	in = in.normalized()
	var c Course
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if err := s.requireUnit(ctx, tx, p, in.UnitID); err != nil {
			return err
		}
		var err error
		if c, err = s.courses.Insert(ctx, tx, in); err != nil {
			return err
		}
		return s.record(ctx, tx, p, ActionCourseCreated, c.ID, map[string]any{"title": c.Title})
	})
	return c, err
}

func (s Service) Get(ctx context.Context, p auth.Principal, id string) (CourseDetails, error) {
	var d CourseDetails
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		c, st, err := s.gate.Visible(ctx, tx, p, id)
		if err != nil {
			return err
		}
		teachers, err := s.roster.Teachers(ctx, tx, id)
		if err != nil {
			return err
		}
		people, err := s.people.ByIDs(ctx, tx, teachers)
		if err != nil {
			return err
		}
		unitIDs, err := s.roster.Units(ctx, tx, id)
		if err != nil {
			return err
		}
		units := make([]structure.Unit, 0, len(unitIDs))
		for _, unitID := range unitIDs {
			u, err := s.units.ByID(ctx, tx, unitID)
			if err != nil {
				return err
			}
			units = append(units, u)
		}
		d = CourseDetails{Course: c, Role: st.Shown(), CanManage: st.Manage, Teachers: people, Units: units}
		return nil
	})
	return d, err
}

func (s Service) Update(ctx context.Context, p auth.Principal, id string, in CourseInput) (Course, error) {
	in = in.normalized()
	var c Course
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		before, st, err := s.gate.Visible(ctx, tx, p, id)
		if err != nil {
			return err
		}
		if !st.Teaches() {
			return access.ErrForbidden
		}
		if !sameUnit(before.UnitID, in.UnitID) {
			if !st.Manage {
				return access.ErrForbidden
			}
			if err := s.requireUnit(ctx, tx, p, in.UnitID); err != nil {
				return err
			}
		}
		if c, err = s.courses.Update(ctx, tx, id, in); err != nil {
			return err
		}
		return s.record(ctx, tx, p, ActionCourseUpdated, id, map[string]any{"before": before, "after": c})
	})
	return c, err
}

func (s Service) SetArchived(ctx context.Context, p auth.Principal, id string, archived bool) error {
	return s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		c, err := s.gate.Managed(ctx, tx, p, id)
		if err != nil {
			return err
		}
		if (c.ArchivedAt != nil) == archived {
			return nil
		}
		if err := s.courses.SetArchived(ctx, tx, id, archived); err != nil {
			return err
		}
		action := ActionCourseRestored
		if archived {
			action = ActionCourseArchived
		}
		return s.record(ctx, tx, p, action, id, map[string]any{"title": c.Title})
	})
}

func (s Service) Change(ctx context.Context, p auth.Principal, list List, courseID, memberID string, add bool) error {
	return s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		c, st, err := s.gate.Visible(ctx, tx, p, courseID)
		if err != nil {
			return err
		}
		allowed := st.Manage || (list == ListStudents && st.Role == RoleTeacher)
		if !allowed {
			return access.ErrForbidden
		}
		if add && c.ArchivedAt != nil {
			return ErrArchived
		}
		var changed bool
		if add {
			changed, err = s.roster.Add(ctx, tx, list, courseID, memberID)
		} else {
			changed, err = s.roster.Remove(ctx, tx, list, courseID, memberID)
		}
		if err != nil || !changed {
			return err
		}
		added, removed := list.actions()
		action := removed
		if add {
			action = added
		}
		return s.record(ctx, tx, p, action, courseID, map[string]any{"title": c.Title, "member": memberID})
	})
}

func (s Service) Students(ctx context.Context, p auth.Principal, courseID string) ([]identity.Person, error) {
	var list []identity.Person
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		_, st, err := s.gate.Visible(ctx, tx, p, courseID)
		if err != nil {
			return err
		}
		if !st.Teaches() {
			return access.ErrForbidden
		}
		ids, err := s.gate.StudentIDs(ctx, tx, courseID)
		if err != nil {
			return err
		}
		list, err = s.people.ByIDs(ctx, tx, ids)
		return err
	})
	return list, err
}

func (s Service) requireUnit(ctx context.Context, tx pgx.Tx, p auth.Principal, unitID *string) error {
	if unitID != nil {
		if _, err := s.units.ByID(ctx, tx, *unitID); err != nil {
			return ErrUnitNotFound
		}
	}
	return s.access.Require(ctx, tx, p, s.gate.Scope(p, unitID), ManageCourses)
}

func (s Service) record(ctx context.Context, tx pgx.Tx, p auth.Principal, action, courseID string, details any) error {
	return s.log.Record(ctx, tx, p.UserID, audit.Change{Action: action, TargetKind: "study.course", TargetID: courseID, Details: details})
}

func sameUnit(a, b *string) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}
