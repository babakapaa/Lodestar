package study

import (
	"context"
	"slices"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
)

func (s AssignmentService) Preview(ctx context.Context, p auth.Principal, courseID, studentID string, at time.Time) ([]AssignmentItem, error) {
	items := []AssignmentItem{}
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		_, st, err := s.gate.Visible(ctx, tx, p, courseID)
		if err != nil {
			return err
		}
		if !st.Teaches() {
			return access.ErrForbidden
		}
		if studentID != "" {
			students, err := s.gate.StudentIDs(ctx, tx, courseID)
			if err != nil {
				return err
			}
			if !slices.Contains(students, studentID) {
				return ErrNotStudent
			}
		}
		all, err := s.assignments.ByCourse(ctx, tx, courseID, false)
		if err != nil {
			return err
		}
		var list []Assignment
		var ids []string
		for _, a := range all {
			if a.PublishedAt != nil && !a.PublishedAt.After(at) {
				list = append(list, a)
				ids = append(ids, a.ID)
			}
		}
		items = make([]AssignmentItem, len(list))
		if studentID == "" {
			for i, a := range list {
				items[i] = AssignmentItem{Assignment: a, Mine: &Progress{}}
			}
			return nil
		}
		if err := s.fillProgress(ctx, tx, studentID, ids, list, items); err != nil {
			return err
		}
		for _, item := range items {
			item.Mine.rewind(at)
		}
		return nil
	})
	return items, err
}
