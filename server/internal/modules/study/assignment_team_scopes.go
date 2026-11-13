package study

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/module"
)

const TeamScopeAssignment = "study.assignment"

type AssignmentTeamScopes struct {
	assignments AssignmentStore
	courses     CourseStore
	rosters     CourseRosters
}

func (AssignmentTeamScopes) Kind() string {
	return TeamScopeAssignment
}

func (s AssignmentTeamScopes) TeamScope(ctx context.Context, tx pgx.Tx, id string) (module.TeamScope, error) {
	a, err := s.assignments.ByID(ctx, tx, id, false)
	if errors.Is(err, ErrAssignmentNotFound) {
		return module.TeamScope{}, module.ErrTeamScopeNotFound
	}
	if err != nil {
		return module.TeamScope{}, err
	}
	if a.Team == nil || !a.Published() {
		return module.TeamScope{}, module.ErrTeamScopeNotFound
	}
	c, err := s.courses.ByID(ctx, tx, a.CourseID)
	if err != nil {
		return module.TeamScope{}, err
	}
	roster, err := s.rosters.One(ctx, tx, a.CourseID)
	if err != nil {
		return module.TeamScope{}, err
	}
	return module.TeamScope{
		Kind: TeamScopeAssignment, ID: a.ID, Title: a.Title, Context: c.Title, Link: assignmentLink(a.ID),
		Leads: roster.Leads, Members: roster.Members,
		MinSize: a.Team.MinSize, MaxSize: a.Team.MaxSize, SelfFormed: a.Team.SelfFormed,
		FreeLeaveUntil: a.Team.FreeLeaveUntil, DueAt: a.DueAt, Archived: c.ArchivedAt != nil,
	}, nil
}
