package study

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/module"
)

const BriefCourse = "course"

type CourseBriefer struct {
	gate CourseGate
}

func (CourseBriefer) Kind() string {
	return BriefCourse
}

func (b CourseBriefer) Brief(ctx context.Context, tx pgx.Tx, p auth.Principal, id string) (module.Brief, error) {
	c, st, err := b.gate.Visible(ctx, tx, p, id)
	if errors.Is(err, ErrCourseNotFound) {
		return module.Brief{}, module.ErrBriefNotFound
	}
	if err != nil {
		return module.Brief{}, err
	}
	return module.Brief{
		Kind: BriefCourse, ID: c.ID, Title: c.Title, Description: c.Description,
		Course: c.Title, CourseID: c.ID, Teaches: st.Teaches(), Mode: c.NavigatorMode,
	}, nil
}
