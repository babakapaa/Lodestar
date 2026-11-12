package study_test

import (
	"errors"
	"testing"
	"time"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/modules/study"
)

func TestScheduleFollowsGroupsAndChanges(t *testing.T) {
	w := newWorld(t)
	stream := w.unit(nil, "stream", "ИС-2023")
	groupA := w.unit(&stream, "group", "ИС-21")
	groupB := w.unit(&stream, "group", "ИС-22")
	alice := w.user("alice", "alice@kaznu.kz")
	bob := w.user("bob", "bob@kaznu.kz")
	teacher := w.user("dana", "dana@kaznu.kz")
	w.join(groupA, alice)
	w.join(groupB, bob)
	c := w.course(nil, "Сети")
	w.add(study.ListUnits, c, stream.ID)
	w.add(study.ListTeachers, c, teacher.UserID)
	schedule := study.NewScheduleService(w.kernel())

	lecture, err := schedule.Create(w.ctx, w.admin, c.ID, study.SeriesInput{
		Kind: "lecture", Room: "301", StartsOn: day(t, "2026-09-01"), EndsOn: day(t, "2026-12-31"),
		Weekday: 2, StartTime: 9 * 60, Minutes: 80,
	})
	if err != nil {
		t.Fatalf("lecture: %v", err)
	}
	if _, err := schedule.Create(w.ctx, teacher, c.ID, study.SeriesInput{
		UnitID: &groupA.ID, Kind: "lab", StartsOn: day(t, "2026-09-01"), EndsOn: day(t, "2026-12-31"),
		Weekday: 4, StartTime: 11 * 60, Minutes: 80,
	}); err != nil {
		t.Fatalf("teacher adds lab: %v", err)
	}
	if _, err := schedule.Create(w.ctx, alice, c.ID, study.SeriesInput{
		Kind: "lab", StartsOn: day(t, "2026-09-01"), EndsOn: day(t, "2026-12-31"), Weekday: 5, Minutes: 80,
	}); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("student edits schedule: err = %v", err)
	}

	week := func(p study.Day) (study.Day, study.Day) { return p, p.AddDays(6) }
	from, to := week(day(t, "2026-09-07"))
	aliceWeek, err := schedule.Mine(w.ctx, alice, from, to)
	if err != nil || len(aliceWeek) != 2 {
		t.Fatalf("alice week = %+v, %v", aliceWeek, err)
	}
	bobWeek, err := schedule.Mine(w.ctx, bob, from, to)
	if err != nil || len(bobWeek) != 1 || bobWeek[0].Kind != "lecture" {
		t.Fatalf("bob week = %+v, %v", bobWeek, err)
	}
	teacherWeek, _ := schedule.Mine(w.ctx, teacher, from, to)
	if len(teacherWeek) != 2 {
		t.Fatalf("teacher week = %+v", teacherWeek)
	}

	if _, err := schedule.PutChange(w.ctx, teacher, lecture.ID, day(t, "2026-09-09"), study.ChangeInput{Status: study.ChangeCancelled}); !errors.Is(err, study.ErrInvalidChange) {
		t.Fatalf("change on a day without lesson: err = %v", err)
	}
	if _, err := schedule.PutChange(w.ctx, teacher, lecture.ID, day(t, "2026-09-08"), study.ChangeInput{Status: study.ChangeCancelled, Note: "Конференция"}); err != nil {
		t.Fatalf("cancel: %v", err)
	}
	changed := w.sent.of(study.NotifyScheduleChanged)
	last := changed[len(changed)-1]
	if len(last.UserIDs) != 2 || last.Params["date"] != "2026-09-08" {
		t.Fatalf("change notice = %+v", last)
	}
	for _, id := range last.UserIDs {
		if id == teacher.UserID {
			t.Fatal("actor told about own change")
		}
	}
	bobWeek, _ = schedule.Mine(w.ctx, bob, from, to)
	if bobWeek[0].Status != study.LessonCancelled || bobWeek[0].Note != "Конференция" {
		t.Fatalf("bob sees = %+v", bobWeek[0])
	}
	if err := schedule.RevertChange(w.ctx, teacher, lecture.ID, day(t, "2026-09-08")); err != nil {
		t.Fatalf("revert: %v", err)
	}
	bobWeek, _ = schedule.Mine(w.ctx, bob, from, to)
	if bobWeek[0].Status != study.LessonPlanned {
		t.Fatalf("after revert = %+v", bobWeek[0])
	}

	w.settings.values[study.TeachersEditSchedule] = false
	if _, err := schedule.PutChange(w.ctx, teacher, lecture.ID, day(t, "2026-09-08"), study.ChangeInput{Status: study.ChangeCancelled}); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("teacher edits when university forbids: err = %v", err)
	}
	if _, err := schedule.Mine(w.ctx, alice, day(t, "2026-09-01"), day(t, "2026-12-31")); !errors.Is(err, study.ErrInvalidRange) {
		t.Fatalf("too long range: err = %v", err)
	}
	if _, err := schedule.Series(w.ctx, w.user("x", "x@kaznu.kz"), c.ID); !errors.Is(err, study.ErrCourseNotFound) {
		t.Fatalf("outsider series: err = %v", err)
	}
}

func TestAgendaAndReminders(t *testing.T) {
	w := newWorld(t)
	teacher := w.user("dana", "dana@kaznu.kz")
	student := w.user("aigerim", "aigerim@kaznu.kz")
	c := w.course(nil, "Химия")
	w.add(study.ListTeachers, c, teacher.UserID)
	w.add(study.ListStudents, c, student.UserID)
	soon := time.Now().Add(5 * time.Hour)
	later := time.Now().Add(72 * time.Hour)
	first, err := w.work.Create(w.ctx, teacher, c.ID, study.AssignmentInput{Title: "Отчёт", DueAt: &soon, Published: true})
	if err != nil {
		t.Fatalf("first: %v", err)
	}
	if _, err := w.work.Create(w.ctx, teacher, c.ID, study.AssignmentInput{Title: "Эссе", DueAt: &later, Published: true}); err != nil {
		t.Fatalf("second: %v", err)
	}
	schedule := study.NewScheduleService(w.kernel())

	agenda, err := schedule.Agenda(w.ctx, student, nil)
	if err != nil || len(agenda.Deadlines) != 2 || agenda.TimeZone != "Asia/Almaty" || agenda.Lessons == nil {
		t.Fatalf("agenda = %+v, %v", agenda, err)
	}
	if teacherAgenda, _ := schedule.Agenda(w.ctx, teacher, nil); len(teacherAgenda.Deadlines) != 0 {
		t.Fatalf("teacher deadlines = %+v", teacherAgenda.Deadlines)
	}

	reminder := study.NewDeadlineReminder(w.kernel())
	if err := reminder.RemindAll(w.ctx); err != nil {
		t.Fatalf("remind: %v", err)
	}
	if err := reminder.RemindAll(w.ctx); err != nil {
		t.Fatalf("remind again: %v", err)
	}
	reminded := w.sent.of(study.NotifyDeadline)
	if len(reminded) != 1 || reminded[0].UserIDs[0] != student.UserID || reminded[0].Params["assignment"] != "Отчёт" {
		t.Fatalf("reminders = %+v", reminded)
	}
	if _, _, err := w.work.Submit(w.ctx, student, first.ID, study.SubmissionInput{ClientID: "6f1d7a52-3c38-4a8e-9a43-1f0e5f3c1a11", Body: "Готово"}); err != nil {
		t.Fatalf("submit: %v", err)
	}
	agenda, _ = schedule.Agenda(w.ctx, student, nil)
	if len(agenda.Deadlines) != 2 || agenda.Deadlines[0].SubmittedAt == nil {
		t.Fatalf("agenda after submit = %+v", agenda.Deadlines)
	}
}
