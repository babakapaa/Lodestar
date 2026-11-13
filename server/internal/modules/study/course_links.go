package study

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/audit"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/module"
)

const LinkCourse = "study.course"

const syncSource = "lti"

type CourseLinks struct {
	gate        CourseGate
	courses     CourseStore
	roster      RosterStore
	assignments AssignmentStore
	log         audit.Log
}

func (CourseLinks) Kind() string {
	return LinkCourse
}

func (CourseLinks) Linkable(ctx context.Context, tx pgx.Tx, p auth.Principal) ([]module.LinkTarget, error) {
	rows, err := tx.Query(ctx, `
		SELECT c.id::text, c.title FROM study.courses c
		JOIN study.course_teachers t ON t.course_id = c.id AND t.user_id::text = $1
		WHERE c.archived_at IS NULL
		ORDER BY c.title, c.id`, p.UserID)
	if err != nil {
		return nil, fmt.Errorf("linkable courses: %w", err)
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (module.LinkTarget, error) {
		var id, title string
		if err := row.Scan(&id, &title); err != nil {
			return module.LinkTarget{}, err
		}
		return courseTarget(id, title), nil
	})
}

func (l CourseLinks) Describe(ctx context.Context, tx pgx.Tx, id string) (module.LinkTarget, error) {
	c, err := l.courses.ByID(ctx, tx, id)
	if errors.Is(err, ErrCourseNotFound) {
		return module.LinkTarget{}, module.ErrLinkTargetNotFound
	}
	if err != nil {
		return module.LinkTarget{}, err
	}
	return courseTarget(c.ID, c.Title), nil
}

func (l CourseLinks) CanLink(ctx context.Context, tx pgx.Tx, p auth.Principal, id string) error {
	c, st, err := l.gate.Visible(ctx, tx, p, id)
	if errors.Is(err, ErrCourseNotFound) {
		return module.ErrLinkTargetNotFound
	}
	if err != nil {
		return err
	}
	if c.ArchivedAt != nil || (!st.Manage && st.Role != RoleTeacher) {
		return access.ErrForbidden
	}
	return nil
}

func (l CourseLinks) CanManage(ctx context.Context, tx pgx.Tx, p auth.Principal, id string) error {
	_, err := l.gate.Managed(ctx, tx, p, id)
	if errors.Is(err, ErrCourseNotFound) {
		return module.ErrLinkTargetNotFound
	}
	return err
}

func (l CourseLinks) Enroll(ctx context.Context, tx pgx.Tx, actorID, id string, leads, members []string) (int, error) {
	c, err := l.courses.ByID(ctx, tx, id)
	if errors.Is(err, ErrCourseNotFound) {
		return 0, module.ErrLinkTargetNotFound
	}
	if err != nil || c.ArchivedAt != nil {
		return 0, err
	}
	added := 0
	for _, part := range []struct {
		list  List
		users []string
	}{{ListTeachers, leads}, {ListStudents, members}} {
		action, _ := part.list.actions()
		for _, userID := range part.users {
			changed, err := l.roster.Add(ctx, tx, part.list, id, userID)
			if err != nil {
				return added, err
			}
			if !changed {
				continue
			}
			added++
			err = l.log.Record(ctx, tx, actorID, audit.Change{
				Action: action, TargetKind: "study.course", TargetID: id,
				Details: map[string]any{"title": c.Title, "member": userID, "source": syncSource},
			})
			if err != nil {
				return added, err
			}
		}
	}
	return added, nil
}

func (l CourseLinks) PutWork(ctx context.Context, tx pgx.Tx, actorID, id, workID string, w module.ExternalWork) (module.WorkRef, error) {
	if workID != "" {
		a, err := l.assignments.Retitle(ctx, tx, id, workID, w)
		if !errors.Is(err, ErrAssignmentNotFound) {
			return module.WorkRef{Kind: TeamScopeAssignment, ID: a.ID}, err
		}
	}
	in, err := AssignmentInput{Title: w.Title, DueAt: w.DueAt, Grading: GradingPoints, MaxPoints: &w.Max, Published: true}.normalized()
	if err != nil {
		return module.WorkRef{}, err
	}
	a, err := l.assignments.Insert(ctx, tx, id, actorID, in)
	if err != nil {
		return module.WorkRef{}, err
	}
	return module.WorkRef{Kind: TeamScopeAssignment, ID: a.ID}, l.log.Record(ctx, tx, actorID, audit.Change{
		Action: ActionAssignmentCreated, TargetKind: "study.assignment", TargetID: a.ID,
		Details: map[string]any{"title": a.Title, "course": id, "source": syncSource},
	})
}

func courseTarget(id, title string) module.LinkTarget {
	return module.LinkTarget{Kind: LinkCourse, ID: id, Title: title, Path: "/study/courses/" + id}
}
