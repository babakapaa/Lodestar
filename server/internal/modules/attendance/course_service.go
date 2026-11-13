package attendance

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/module"
)

const (
	dayLayout    = "2006-01-02"
	defaultSpan  = 120
	longestRange = 370
)

func (s Service) Course(ctx context.Context, p auth.Principal, courseID, from, to string) (CourseAttendance, error) {
	var out CourseAttendance
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		var err error
		out, _, err = s.course(ctx, tx, p, courseID, from, to)
		return err
	})
	return out, err
}

func (s Service) course(ctx context.Context, tx pgx.Tx, p auth.Principal, courseID, from, to string) (CourseAttendance, map[string]map[string]*MarkStatus, error) {
	from, to, err := s.span(ctx, tx, from, to)
	if err != nil {
		return CourseAttendance{}, nil, err
	}
	if s.lessons == nil {
		return CourseAttendance{}, nil, ErrLessonNotFound
	}
	c, err := s.lessons.Course(ctx, tx, courseID, from, to)
	if errors.Is(err, module.ErrLessonNotFound) {
		return CourseAttendance{}, nil, ErrLessonNotFound
	}
	if err != nil {
		return CourseAttendance{}, nil, err
	}
	if err := s.lead(ctx, tx, p, c.Leads); err != nil {
		if errors.Is(err, access.ErrForbidden) {
			return CourseAttendance{}, nil, ErrLessonNotFound
		}
		return CourseAttendance{}, nil, err
	}
	alert, err := s.integer(ctx, tx, AbsentAlert)
	if err != nil {
		return CourseAttendance{}, nil, err
	}
	sessions, err := s.sessions.OfCourse(ctx, tx, courseID, from, to)
	if err != nil {
		return CourseAttendance{}, nil, err
	}
	bySeries := make(map[string]sessionRow, len(sessions))
	ids := make([]string, len(sessions))
	for i, r := range sessions {
		bySeries[r.SeriesID+"/"+r.Date] = r
		ids[i] = r.ID
	}
	marks, err := s.marks.OfSessions(ctx, tx, ids)
	if err != nil {
		return CourseAttendance{}, nil, err
	}
	grid := map[string]map[string]*MarkStatus{}
	for _, m := range marks {
		if grid[m.SessionID] == nil {
			grid[m.SessionID] = map[string]*MarkStatus{}
		}
		grid[m.SessionID][m.UserID] = m.Status
	}
	out := CourseAttendance{CourseID: c.ID, CourseTitle: c.Title, From: from, To: to, AbsentAlert: alert, Lessons: []LessonSummary{}, Students: []StudentRecord{}}
	byStudent := map[string]*StudentRecord{}
	people := append([]string{}, c.Members...)
	for _, m := range marks {
		people = append(people, m.UserID)
	}
	names, err := s.names(ctx, tx, people)
	if err != nil {
		return out, nil, err
	}
	student := func(id string) *StudentRecord {
		if r, ok := byStudent[id]; ok {
			return r
		}
		byStudent[id] = &StudentRecord{UserID: id, Name: names[id]}
		return byStudent[id]
	}
	for _, id := range c.Members {
		student(id)
	}
	for _, l := range c.Lessons {
		summary := LessonSummary{SeriesID: l.SeriesID, Date: l.Date, StartsAt: l.StartsAt, Room: l.Room, Cancelled: l.Cancelled}
		r, ok := bySeries[l.SeriesID+"/"+l.Date]
		if ok {
			summary.SessionID, summary.Open = &r.ID, r.open()
			for userID, status := range grid[r.ID] {
				summary.Counts.add(status)
				rec := student(userID)
				rec.Counts.add(status)
				if r.open() || status == nil {
					continue
				}
				switch *status {
				case Absent:
					rec.Streak++
				case Present, Late:
					rec.Streak = 0
				}
			}
		}
		out.Lessons = append(out.Lessons, summary)
	}
	for _, rec := range byStudent {
		held := rec.Counts.Present + rec.Counts.Late + rec.Counts.Absent
		if held > 0 {
			rate := float64(rec.Counts.Present+rec.Counts.Late) / float64(held)
			rec.Rate = &rate
		}
		rec.LongAbsent = rec.Streak >= alert
		out.Students = append(out.Students, *rec)
	}
	sort.Slice(out.Students, func(i, j int) bool {
		if out.Students[i].Name != out.Students[j].Name {
			return out.Students[i].Name < out.Students[j].Name
		}
		return out.Students[i].UserID < out.Students[j].UserID
	})
	return out, grid, nil
}

func (s Service) span(ctx context.Context, tx pgx.Tx, from, to string) (string, string, error) {
	today := s.now().In(s.location(ctx, tx)).Format(dayLayout)
	if to == "" {
		to = today
	}
	last, err := time.Parse(dayLayout, to)
	if err != nil {
		return "", "", ErrInvalidRange
	}
	if from == "" {
		from = last.AddDate(0, 0, -defaultSpan).Format(dayLayout)
	}
	first, err := time.Parse(dayLayout, from)
	if err != nil || last.Before(first) || last.Sub(first) > longestRange*24*time.Hour {
		return "", "", ErrInvalidRange
	}
	return from, to, nil
}
