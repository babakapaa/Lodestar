package study

import (
	"sort"
	"time"
)

type Calendar struct {
	Location *time.Location
}

func (c Calendar) Expand(series []Series, changes []Change, titles map[string]string, from, to Day) []Lesson {
	byKey := make(map[string]Change, len(changes))
	for _, ch := range changes {
		byKey[ch.SeriesID+"/"+ch.Date.String()] = ch
	}
	var lessons []Lesson
	for _, s := range series {
		for _, day := range c.days(s, from, to) {
			lesson := Lesson{
				SeriesID: s.ID, CourseID: s.CourseID, CourseTitle: titles[s.CourseID], UnitID: s.UnitID,
				Kind: s.Kind, Date: day, Room: s.Room, TeacherID: s.TeacherID, Status: LessonPlanned,
			}
			lesson.StartsAt = day.At(int(s.StartTime), c.Location)
			minutes := s.Minutes
			if ch, ok := byKey[s.ID+"/"+day.String()]; ok {
				lesson.Note = ch.Note
				if ch.Status == ChangeCancelled {
					lesson.Status = LessonCancelled
				} else {
					lesson.Status = LessonChanged
					if ch.StartsAt != nil {
						lesson.StartsAt = ch.StartsAt.In(c.Location)
					}
					if ch.Minutes != nil {
						minutes = *ch.Minutes
					}
					if ch.Room != nil {
						lesson.Room = *ch.Room
					}
					if ch.TeacherID != nil {
						lesson.TeacherID = ch.TeacherID
					}
				}
			}
			lesson.EndsAt = lesson.StartsAt.Add(time.Duration(minutes) * time.Minute)
			lessons = append(lessons, lesson)
		}
	}
	sort.SliceStable(lessons, func(i, j int) bool {
		return lessons[i].StartsAt.Before(lessons[j].StartsAt)
	})
	return lessons
}

func (c Calendar) Occurs(s Series, day Day) bool {
	days := c.days(s, day, day)
	return len(days) == 1
}

func (Calendar) days(s Series, from, to Day) []Day {
	first := s.StartsOn.AddDays((s.Weekday - s.StartsOn.Weekday() + 7) % 7)
	step := 7 * s.EveryWeeks
	start := first
	if start.Before(from.Time) {
		gap := int(from.Sub(first.Time).Hours() / 24)
		start = first.AddDays((gap + step - 1) / step * step)
	}
	end := to
	if s.EndsOn.Before(end.Time) {
		end = s.EndsOn
	}
	var out []Day
	for d := start; !d.After(end.Time); d = d.AddDays(step) {
		out = append(out, d)
	}
	return out
}
