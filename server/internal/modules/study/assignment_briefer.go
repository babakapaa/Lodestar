package study

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/module"
)

const BriefAssignment = "assignment"

type AssignmentBriefer struct {
	service AssignmentService
}

func (AssignmentBriefer) Kind() string {
	return BriefAssignment
}

func (b AssignmentBriefer) Brief(ctx context.Context, tx pgx.Tx, p auth.Principal, id string) (module.Brief, error) {
	a, c, st, err := b.service.visible(ctx, tx, p, id, false)
	if errors.Is(err, ErrAssignmentNotFound) {
		return module.Brief{}, module.ErrBriefNotFound
	}
	if err != nil {
		return module.Brief{}, err
	}
	mode := c.NavigatorMode
	if a.NavigatorMode != nil {
		mode = *a.NavigatorMode
	}
	teachers, err := RosterStore{}.Teachers(ctx, tx, c.ID)
	if err != nil {
		return module.Brief{}, err
	}
	return module.Brief{
		Kind: BriefAssignment, ID: a.ID, Title: a.Title, Description: a.Description,
		Course: c.Title, CourseID: c.ID, Teaches: st.Teaches(), DueAt: a.DueAt, Mode: mode, Teachers: teachers,
	}, nil
}
