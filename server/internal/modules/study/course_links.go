package study

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/module"
)

const LinkCourse = "study.course"

type CourseLinks struct {
	gate    CourseGate
	courses CourseStore
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

func courseTarget(id, title string) module.LinkTarget {
	return module.LinkTarget{Kind: LinkCourse, ID: id, Title: title, Path: "/study/courses/" + id}
}
