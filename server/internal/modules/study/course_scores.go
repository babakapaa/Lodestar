package study

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/module"
)

type CourseScores struct {
	CourseRosters
	assignments AssignmentStore
	results     ResultStore
}

func (s CourseScores) ScoredWorks(ctx context.Context, tx pgx.Tx, courseID string) ([]module.ScoredWork, error) {
	list, err := s.assignments.ByCourse(ctx, tx, courseID, false)
	if err != nil {
		return nil, err
	}
	results, err := s.results.ByCourse(ctx, tx, courseID)
	if err != nil {
		return nil, err
	}
	works := make([]module.ScoredWork, 0, len(list))
	for _, a := range list {
		if a.Grading == GradingNone {
			continue
		}
		works = append(works, module.ScoredWork{
			Kind: TeamScopeAssignment, ID: a.ID, Title: a.Title,
			Project: a.Team != nil || a.Defense != nil, Max: a.maxScore(), Results: results[a.ID],
		})
	}
	return works, nil
}
