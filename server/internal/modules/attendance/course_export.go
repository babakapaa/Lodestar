package attendance

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
)

type CourseTable struct {
	Title  string
	Header []string
	Rows   [][]any
}

func (s Service) ExportCourse(ctx context.Context, p auth.Principal, courseID, from, to string) (CourseTable, error) {
	var out CourseTable
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		c, grid, err := s.course(ctx, tx, p, courseID, from, to)
		if err != nil {
			return err
		}
		loc := s.location(ctx, tx)
		out.Title = c.CourseTitle
		out.Header = []string{"student"}
		var held []LessonSummary
		for _, l := range c.Lessons {
			if l.SessionID != nil {
				held = append(held, l)
				out.Header = append(out.Header, l.StartsAt.In(loc).Format("2006-01-02 15:04"))
			}
		}
		out.Header = append(out.Header, "present", "late", "absent", "excused", "rate")
		for _, st := range c.Students {
			row := []any{st.Name}
			for _, l := range held {
				var cell any
				if status := grid[*l.SessionID][st.UserID]; status != nil {
					cell = string(*status)
				}
				row = append(row, cell)
			}
			var rate any
			if st.Rate != nil {
				rate = *st.Rate
			}
			out.Rows = append(out.Rows, append(row, st.Counts.Present, st.Counts.Late, st.Counts.Absent, st.Counts.Excused, rate))
		}
		return nil
	})
	return out, err
}
