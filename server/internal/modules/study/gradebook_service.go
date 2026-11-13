package study

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

type GradebookService struct {
	tenants     database.TenantTx
	gate        CourseGate
	assignments AssignmentStore
	results     ResultStore
	people      identity.People
}

func NewGradebookService(tenants database.TenantTx, acc access.Service) GradebookService {
	return GradebookService{tenants: tenants, gate: CourseGate{access: acc}}
}

func (s GradebookService) Build(ctx context.Context, p auth.Principal, courseID string) (Gradebook, error) {
	var book Gradebook
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		_, st, err := s.gate.Visible(ctx, tx, p, courseID)
		if err != nil {
			return err
		}
		if !st.Teaches() {
			return access.ErrForbidden
		}
		list, err := s.assignments.ByCourse(ctx, tx, courseID, false)
		if err != nil {
			return err
		}
		results, err := s.results.ByCourse(ctx, tx, courseID)
		if err != nil {
			return err
		}
		ids, err := s.gate.StudentIDs(ctx, tx, courseID)
		if err != nil {
			return err
		}
		students, err := s.people.ByIDs(ctx, tx, ids)
		if err != nil {
			return err
		}
		book = s.table(list, results, students)
		return nil
	})
	return book, err
}

func (s GradebookService) table(list []Assignment, results map[string][]module.WorkResult, students []identity.Person) Gradebook {
	book := Gradebook{Header: []string{"name", "email"}}
	var graded []Assignment
	for _, a := range list {
		if a.Grading == GradingNone {
			continue
		}
		graded = append(graded, a)
		book.Header = append(book.Header, fmt.Sprintf("%s (%g)", a.Title, a.maxScore()))
	}
	scores := map[string]map[string]float64{}
	for assignmentID, rs := range results {
		for _, r := range rs {
			if scores[r.UserID] == nil {
				scores[r.UserID] = map[string]float64{}
			}
			scores[r.UserID][assignmentID] = r.Score
		}
	}
	for _, person := range students {
		row := []any{person.DisplayName, person.Email}
		for _, a := range graded {
			if v, ok := scores[person.UserID][a.ID]; ok {
				row = append(row, v)
			} else {
				row = append(row, nil)
			}
		}
		book.Rows = append(book.Rows, row)
	}
	return book
}
