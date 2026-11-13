package study_test

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/modules/study"
)

func TestLessonSourceKnowsGroupsAndCancellations(t *testing.T) {
	w := newWorld(t)
	stream := w.unit(nil, "stream", "ИС-2023")
	group := w.unit(&stream, "group", "ИС-21")
	other := w.unit(&stream, "group", "ИС-22")
	alice := w.user("alice", "alice@kaznu.kz")
	bob := w.user("bob", "bob@kaznu.kz")
	teacher := w.user("dana", "dana@kaznu.kz")
	w.join(group, alice)
	w.join(other, bob)
	c := w.course(nil, "Сети")
	w.add(study.ListUnits, c, stream.ID)
	w.add(study.ListTeachers, c, teacher.UserID)
	schedule := study.NewScheduleService(w.kernel())
	lecture, err := schedule.Create(w.ctx, w.admin, c.ID, study.SeriesInput{
		Kind: "lecture", Room: "301", StartsOn: day(t, "2026-09-01"), EndsOn: day(t, "2026-12-31"), Weekday: 2, StartTime: 9 * 60, Minutes: 80,
	})
	if err != nil {
		t.Fatalf("lecture: %v", err)
	}
	lab, err := schedule.Create(w.ctx, w.admin, c.ID, study.SeriesInput{
		UnitID: &group.ID, Kind: "lab", StartsOn: day(t, "2026-09-01"), EndsOn: day(t, "2026-12-31"), Weekday: 4, StartTime: 11 * 60, Minutes: 80,
	})
	if err != nil {
		t.Fatalf("lab: %v", err)
	}
	if _, err := schedule.PutChange(w.ctx, w.admin, lecture.ID, day(t, "2026-09-15"), study.ChangeInput{Status: study.ChangeCancelled}); err != nil {
		t.Fatalf("cancel: %v", err)
	}
	source := study.Module{}.LessonSource(w.kernel())
	var labDay module.LessonRef
	var missing error
	var course module.LessonCourse
	err = w.tx.Run(w.ctx, w.tenant, func(tx pgx.Tx) error {
		var err error
		if labDay, err = source.Lesson(w.ctx, tx, lab.ID, "2026-09-10"); err != nil {
			return err
		}
		_, missing = source.Lesson(w.ctx, tx, lab.ID, "2026-09-11")
		course, err = source.Course(w.ctx, tx, c.ID, "2026-09-07", "2026-09-20")
		return err
	})
	if err != nil {
		t.Fatalf("lessons: %v", err)
	}
	if labDay.Room != "" || len(labDay.Members) != 1 || labDay.Members[0] != alice.UserID || labDay.Leads[0] != teacher.UserID || labDay.CourseTitle != "Сети" {
		t.Fatalf("lab = %+v", labDay)
	}
	if !errors.Is(missing, module.ErrLessonNotFound) {
		t.Fatalf("no lesson on friday: %v", missing)
	}
	if len(course.Lessons) != 4 || len(course.Members) != 2 || !course.Lessons[2].Cancelled || len(course.Lessons[0].Members) != 2 {
		t.Fatalf("course = %+v", course)
	}
}
