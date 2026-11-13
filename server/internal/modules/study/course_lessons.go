package study

import (
	"context"
	"encoding/json"
	"errors"
	"slices"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/core/structure"
)

type CourseLessons struct {
	settings  module.Settings
	series    SeriesStore
	changes   ChangeStore
	rosters   CourseRosters
	directory structure.Directory
}

func (c CourseLessons) Lesson(ctx context.Context, tx pgx.Tx, seriesID, date string) (module.LessonRef, error) {
	day, err := ParseDay(date)
	if err != nil {
		return module.LessonRef{}, module.ErrLessonNotFound
	}
	s, err := c.series.ByID(ctx, tx, seriesID)
	if errors.Is(err, ErrSeriesNotFound) {
		return module.LessonRef{}, module.ErrLessonNotFound
	}
	if err != nil {
		return module.LessonRef{}, err
	}
	if !(Calendar{}).Occurs(s, day) {
		return module.LessonRef{}, module.ErrLessonNotFound
	}
	course, err := c.Course(ctx, tx, s.CourseID, date, date)
	if err != nil {
		return module.LessonRef{}, err
	}
	for _, l := range course.Lessons {
		if l.SeriesID == seriesID {
			return l, nil
		}
	}
	return module.LessonRef{}, module.ErrLessonNotFound
}

func (c CourseLessons) Course(ctx context.Context, tx pgx.Tx, courseID, from, to string) (module.LessonCourse, error) {
	first, err := ParseDay(from)
	if err != nil {
		return module.LessonCourse{}, err
	}
	last, err := ParseDay(to)
	if err != nil {
		return module.LessonCourse{}, err
	}
	roster, err := c.rosters.One(ctx, tx, courseID)
	if errors.Is(err, module.ErrRosterNotFound) {
		return module.LessonCourse{}, module.ErrLessonNotFound
	}
	if err != nil {
		return module.LessonCourse{}, err
	}
	out := module.LessonCourse{ID: roster.ID, Title: roster.Title, Leads: roster.Leads, Members: roster.Members}
	series, err := c.series.ByCourses(ctx, tx, []string{courseID}, first, last)
	if err != nil {
		return out, err
	}
	ids := make([]string, len(series))
	var units []string
	for i, s := range series {
		ids[i] = s.ID
		if s.UnitID != nil {
			units = append(units, *s.UnitID)
		}
	}
	changes, err := c.changes.InRange(ctx, tx, ids, first, last)
	if err != nil {
		return out, err
	}
	under, err := c.directory.UsersUnderEach(ctx, tx, units)
	if err != nil {
		return out, err
	}
	bySeries := make(map[string]Series, len(series))
	for _, s := range series {
		bySeries[s.ID] = s
	}
	lessons := Calendar{Location: c.location(ctx, tx)}.Expand(series, changes, map[string]string{courseID: roster.Title}, first, last)
	for _, l := range lessons {
		ref := module.LessonRef{
			SeriesID: l.SeriesID, CourseID: l.CourseID, CourseTitle: l.CourseTitle, Date: l.Date.String(),
			StartsAt: l.StartsAt, EndsAt: l.EndsAt, Room: l.Room, Cancelled: l.Status == LessonCancelled,
			Leads: slices.Clone(roster.Leads), Members: roster.Members,
		}
		if l.TeacherID != nil && !slices.Contains(ref.Leads, *l.TeacherID) {
			ref.Leads = append(ref.Leads, *l.TeacherID)
		}
		if unit := bySeries[l.SeriesID].UnitID; unit != nil {
			inside := under[*unit]
			ref.Members = slices.DeleteFunc(slices.Clone(roster.Members), func(id string) bool { return !slices.Contains(inside, id) })
		}
		out.Lessons = append(out.Lessons, ref)
	}
	return out, nil
}

func (c CourseLessons) location(ctx context.Context, tx pgx.Tx) *time.Location {
	raw, err := c.settings.Value(ctx, tx, TimeZoneSetting)
	if err != nil {
		return time.UTC
	}
	var name string
	if json.Unmarshal(raw, &name) != nil || name == "" {
		return time.UTC
	}
	loc, err := time.LoadLocation(name)
	if err != nil {
		return time.UTC
	}
	return loc
}
