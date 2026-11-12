package study

import (
	"context"
	"encoding/json"
	"slices"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/audit"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/core/structure"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

const (
	maxRangeDays  = 62
	deadlineAhead = 8
	deadlineBack  = 14
)

type ScheduleService struct {
	tenants     database.TenantTx
	gate        CourseGate
	settings    module.Settings
	notifier    module.Notifier
	courses     CourseStore
	series      SeriesStore
	changes     ChangeStore
	assignments AssignmentStore
	submissions SubmissionStore
	directory   structure.Directory
	log         audit.Log
	now         func() time.Time
}

func NewScheduleService(k module.Kernel) ScheduleService {
	return ScheduleService{
		tenants: k.Tenants, gate: CourseGate{access: k.Access}, settings: k.Settings, notifier: k.Notifier, now: time.Now,
	}
}

func (s ScheduleService) Series(ctx context.Context, p auth.Principal, courseID string) ([]Series, error) {
	var list []Series
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if _, _, err := s.gate.Visible(ctx, tx, p, courseID); err != nil {
			return err
		}
		var err error
		list, err = s.series.OfCourse(ctx, tx, courseID)
		return err
	})
	return list, err
}

func (s ScheduleService) Create(ctx context.Context, p auth.Principal, courseID string, in SeriesInput) (Series, error) {
	in, err := in.normalized()
	if err != nil {
		return Series{}, err
	}
	var created Series
	err = s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		c, err := s.editable(ctx, tx, p, courseID)
		if err != nil {
			return err
		}
		if created, err = s.series.Insert(ctx, tx, courseID, in); err != nil {
			return err
		}
		if err := s.record(ctx, tx, p, ActionSeriesCreated, created.ID, created); err != nil {
			return err
		}
		return s.announce(ctx, tx, p, c, created, created.StartsOn)
	})
	return created, err
}

func (s ScheduleService) Update(ctx context.Context, p auth.Principal, id string, in SeriesInput) (Series, error) {
	in, err := in.normalized()
	if err != nil {
		return Series{}, err
	}
	var updated Series
	err = s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		before, c, err := s.editableSeries(ctx, tx, p, id)
		if err != nil {
			return err
		}
		if updated, err = s.series.Update(ctx, tx, id, in); err != nil {
			return err
		}
		if err := s.record(ctx, tx, p, ActionSeriesUpdated, id, updated); err != nil {
			return err
		}
		return s.announceBoth(ctx, tx, p, c, before, updated)
	})
	return updated, err
}

func (s ScheduleService) Delete(ctx context.Context, p auth.Principal, id string) error {
	return s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		before, c, err := s.editableSeries(ctx, tx, p, id)
		if err != nil {
			return err
		}
		if err := s.announce(ctx, tx, p, c, before, s.today(ctx, tx)); err != nil {
			return err
		}
		if err := s.series.Delete(ctx, tx, id); err != nil {
			return err
		}
		return s.record(ctx, tx, p, ActionSeriesDeleted, id, before)
	})
}

func (s ScheduleService) PutChange(ctx context.Context, p auth.Principal, seriesID string, date Day, in ChangeInput) (Change, error) {
	in, err := in.normalized()
	if err != nil {
		return Change{}, err
	}
	var ch Change
	err = s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		series, c, err := s.editableSeries(ctx, tx, p, seriesID)
		if err != nil {
			return err
		}
		if !(Calendar{}).Occurs(series, date) {
			return ErrInvalidChange
		}
		if ch, err = s.changes.Put(ctx, tx, seriesID, date, in); err != nil {
			return err
		}
		if err := s.record(ctx, tx, p, ActionLessonChanged, seriesID, map[string]any{"date": date, "status": in.Status}); err != nil {
			return err
		}
		return s.announce(ctx, tx, p, c, series, date)
	})
	return ch, err
}

func (s ScheduleService) RevertChange(ctx context.Context, p auth.Principal, seriesID string, date Day) error {
	return s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		series, c, err := s.editableSeries(ctx, tx, p, seriesID)
		if err != nil {
			return err
		}
		removed, err := s.changes.Delete(ctx, tx, seriesID, date)
		if err != nil || !removed {
			return err
		}
		if err := s.record(ctx, tx, p, ActionLessonRestored, seriesID, map[string]any{"date": date}); err != nil {
			return err
		}
		return s.announce(ctx, tx, p, c, series, date)
	})
}

func (s ScheduleService) Mine(ctx context.Context, p auth.Principal, from, to Day) ([]Lesson, error) {
	if err := checkRange(from, to); err != nil {
		return nil, err
	}
	var lessons []Lesson
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		var err error
		lessons, err = s.mine(ctx, tx, p, from, to)
		return err
	})
	return lessons, err
}

func (s ScheduleService) OfCourse(ctx context.Context, p auth.Principal, courseID string, from, to Day) ([]Lesson, error) {
	if err := checkRange(from, to); err != nil {
		return nil, err
	}
	var lessons []Lesson
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		c, st, err := s.gate.Visible(ctx, tx, p, courseID)
		if err != nil {
			return err
		}
		filter, err := s.filter(ctx, tx, p, st.Teaches())
		if err != nil {
			return err
		}
		lessons, err = s.expand(ctx, tx, map[string]string{c.ID: c.Title}, map[string]func(Series) bool{c.ID: filter}, from, to)
		return err
	})
	return lessons, err
}

func (s ScheduleService) Agenda(ctx context.Context, p auth.Principal, date *Day) (Agenda, error) {
	var a Agenda
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		loc, name := s.location(ctx, tx)
		day := DayOf(s.now().In(loc))
		if date != nil {
			day = *date
		}
		lessons, err := s.mine(ctx, tx, p, day, day)
		if err != nil {
			return err
		}
		deadlines, err := s.deadlines(ctx, tx, p, day, loc)
		if err != nil {
			return err
		}
		a = Agenda{Date: day, TimeZone: name, Lessons: orEmpty(lessons), Deadlines: orEmpty(deadlines)}
		return nil
	})
	return a, err
}

func (s ScheduleService) mine(ctx context.Context, tx pgx.Tx, p auth.Principal, from, to Day) ([]Lesson, error) {
	units, err := s.directory.UnitsAbove(ctx, tx, p.UserID)
	if err != nil {
		return nil, err
	}
	courses, err := s.courses.Mine(ctx, tx, p.UserID, units, false)
	if err != nil {
		return nil, err
	}
	titles := make(map[string]string, len(courses))
	filters := make(map[string]func(Series) bool, len(courses))
	studentFilter := unitFilter(units)
	for _, c := range courses {
		titles[c.ID] = c.Title
		filters[c.ID] = studentFilter
		if c.Role == RoleTeacher {
			filters[c.ID] = everything
		}
	}
	return s.expand(ctx, tx, titles, filters, from, to)
}

func (s ScheduleService) expand(ctx context.Context, tx pgx.Tx, titles map[string]string, filters map[string]func(Series) bool, from, to Day) ([]Lesson, error) {
	ids := make([]string, 0, len(titles))
	for id := range titles {
		ids = append(ids, id)
	}
	all, err := s.series.ByCourses(ctx, tx, ids, from, to)
	if err != nil {
		return nil, err
	}
	var series []Series
	var seriesIDs []string
	for _, item := range all {
		if filters[item.CourseID](item) {
			series = append(series, item)
			seriesIDs = append(seriesIDs, item.ID)
		}
	}
	changes, err := s.changes.InRange(ctx, tx, seriesIDs, from, to)
	if err != nil {
		return nil, err
	}
	loc, _ := s.location(ctx, tx)
	return orEmpty(Calendar{Location: loc}.Expand(series, changes, titles, from, to)), nil
}

func (s ScheduleService) filter(ctx context.Context, tx pgx.Tx, p auth.Principal, teaches bool) (func(Series) bool, error) {
	if teaches {
		return everything, nil
	}
	units, err := s.directory.UnitsAbove(ctx, tx, p.UserID)
	if err != nil {
		return nil, err
	}
	return unitFilter(units), nil
}

func (s ScheduleService) deadlines(ctx context.Context, tx pgx.Tx, p auth.Principal, day Day, loc *time.Location) ([]Deadline, error) {
	units, err := s.directory.UnitsAbove(ctx, tx, p.UserID)
	if err != nil {
		return nil, err
	}
	courses, err := s.courses.Mine(ctx, tx, p.UserID, units, false)
	if err != nil {
		return nil, err
	}
	titles := map[string]string{}
	var ids []string
	for _, c := range courses {
		if c.Role == RoleStudent {
			titles[c.ID] = c.Title
			ids = append(ids, c.ID)
		}
	}
	if len(ids) == 0 {
		return nil, nil
	}
	now := s.now()
	due, err := s.assignments.Due(ctx, tx, ids, now.AddDate(0, 0, -deadlineBack), day.AddDays(deadlineAhead).At(0, loc))
	if err != nil {
		return nil, err
	}
	assignmentIDs := make([]string, len(due))
	for i, a := range due {
		assignmentIDs[i] = a.ID
	}
	submitted, err := s.submissions.LastSubmitted(ctx, tx, assignmentIDs, p.UserID)
	if err != nil {
		return nil, err
	}
	var list []Deadline
	for _, a := range due {
		d := Deadline{AssignmentID: a.ID, CourseID: a.CourseID, CourseTitle: titles[a.CourseID], Title: a.Title, DueAt: *a.DueAt}
		if at, ok := submitted[a.ID]; ok {
			d.SubmittedAt = &at
		}
		if d.SubmittedAt != nil && d.DueAt.Before(now) {
			continue
		}
		list = append(list, d)
	}
	return list, nil
}

func (s ScheduleService) editable(ctx context.Context, tx pgx.Tx, p auth.Principal, courseID string) (Course, error) {
	c, st, err := s.gate.Visible(ctx, tx, p, courseID)
	if err != nil {
		return Course{}, err
	}
	err = s.gate.access.Require(ctx, tx, p, s.gate.Scope(p, c.UnitID), ManageSchedule)
	if err == nil {
		return c, nil
	}
	if !isForbidden(err) {
		return Course{}, err
	}
	if st.Role != RoleTeacher {
		return Course{}, access.ErrForbidden
	}
	raw, err := s.settings.Value(ctx, tx, TeachersEditSchedule)
	if err != nil {
		return Course{}, err
	}
	var allowed bool
	if err := json.Unmarshal(raw, &allowed); err != nil || !allowed {
		return Course{}, access.ErrForbidden
	}
	return c, nil
}

func (s ScheduleService) editableSeries(ctx context.Context, tx pgx.Tx, p auth.Principal, id string) (Series, Course, error) {
	series, err := s.series.ByID(ctx, tx, id)
	if err != nil {
		return Series{}, Course{}, err
	}
	c, err := s.editable(ctx, tx, p, series.CourseID)
	if err != nil {
		if isForbidden(err) {
			if _, _, hidden := s.gate.Visible(ctx, tx, p, series.CourseID); hidden != nil {
				return Series{}, Course{}, ErrSeriesNotFound
			}
		}
		return Series{}, Course{}, err
	}
	return series, c, nil
}

func (s ScheduleService) announceBoth(ctx context.Context, tx pgx.Tx, p auth.Principal, c Course, before, after Series) error {
	users, err := s.audience(ctx, tx, p, c, before)
	if err != nil {
		return err
	}
	more, err := s.audience(ctx, tx, p, c, after)
	if err != nil {
		return err
	}
	for _, id := range more {
		if !slices.Contains(users, id) {
			users = append(users, id)
		}
	}
	return s.send(ctx, tx, p, c, users, after.StartsOn)
}

func (s ScheduleService) announce(ctx context.Context, tx pgx.Tx, p auth.Principal, c Course, series Series, date Day) error {
	users, err := s.audience(ctx, tx, p, c, series)
	if err != nil {
		return err
	}
	return s.send(ctx, tx, p, c, users, date)
}

func (s ScheduleService) send(ctx context.Context, tx pgx.Tx, p auth.Principal, c Course, users []string, date Day) error {
	return s.notifier.Notify(ctx, tx, p.TenantID, module.Notification{
		Type: NotifyScheduleChanged, UserIDs: users,
		Params: map[string]any{"course": c.Title, "date": date.String()},
		Link:   "/study/schedule?date=" + date.String(),
	})
}

func (s ScheduleService) audience(ctx context.Context, tx pgx.Tx, p auth.Principal, c Course, series Series) ([]string, error) {
	students, err := s.gate.StudentIDs(ctx, tx, c.ID)
	if err != nil {
		return nil, err
	}
	if series.UnitID != nil {
		inside, err := s.directory.UsersUnder(ctx, tx, []string{*series.UnitID})
		if err != nil {
			return nil, err
		}
		students = slices.DeleteFunc(students, func(id string) bool { return !slices.Contains(inside, id) })
	}
	teachers, err := s.gate.roster.Teachers(ctx, tx, c.ID)
	if err != nil {
		return nil, err
	}
	all := slices.Concat(students, teachers)
	if series.TeacherID != nil {
		all = append(all, *series.TeacherID)
	}
	var out []string
	for _, id := range all {
		if id != p.UserID && !slices.Contains(out, id) {
			out = append(out, id)
		}
	}
	return out, nil
}

func (s ScheduleService) location(ctx context.Context, tx pgx.Tx) (*time.Location, string) {
	raw, err := s.settings.Value(ctx, tx, TimeZoneSetting)
	if err != nil {
		return time.UTC, "UTC"
	}
	var name string
	if json.Unmarshal(raw, &name) != nil || name == "" {
		return time.UTC, "UTC"
	}
	loc, err := time.LoadLocation(name)
	if err != nil {
		return time.UTC, "UTC"
	}
	return loc, name
}

func (s ScheduleService) today(ctx context.Context, tx pgx.Tx) Day {
	loc, _ := s.location(ctx, tx)
	return DayOf(s.now().In(loc))
}

func (s ScheduleService) record(ctx context.Context, tx pgx.Tx, p auth.Principal, action, seriesID string, details any) error {
	return s.log.Record(ctx, tx, p.UserID, audit.Change{Action: action, TargetKind: "study.series", TargetID: seriesID, Details: details})
}

func checkRange(from, to Day) error {
	if to.Before(from.Time) || to.Sub(from.Time).Hours() > maxRangeDays*24 {
		return ErrInvalidRange
	}
	return nil
}

func unitFilter(units []string) func(Series) bool {
	return func(s Series) bool {
		return s.UnitID == nil || slices.Contains(units, *s.UnitID)
	}
}

func everything(Series) bool {
	return true
}

func orEmpty[T any](list []T) []T {
	if list == nil {
		return []T{}
	}
	return list
}
