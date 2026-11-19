package study

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/module"
)

func (s AssignmentService) helpFor(ctx context.Context, tx pgx.Tx, a Assignment, studentID string) (*module.HelpSummary, error) {
	if !a.NavigatorHelp || s.help == nil {
		return nil, nil
	}
	sum, err := s.help.Help(ctx, tx, studentID, BriefAssignment, a.ID)
	if err != nil {
		return nil, err
	}
	return &sum, nil
}
