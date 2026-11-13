package study_test

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/modules/study"
)

func TestAssignmentBriefFollowsCourseMode(t *testing.T) {
	w := newWorld(t)
	teacher, student, outsider := w.user("dana", "dana@kaznu.kz"), w.user("aigerim", "aigerim@kaznu.kz"), w.user("erlan", "erlan@kaznu.kz")
	c, err := w.study.Create(w.ctx, w.admin, study.CourseInput{Title: "Программирование", NavigatorMode: module.ModeProject})
	if err != nil || c.NavigatorMode != module.ModeProject {
		t.Fatalf("course = %+v, %v", c, err)
	}
	if _, err := w.study.Create(w.ctx, w.admin, study.CourseInput{Title: "x", NavigatorMode: "exam"}); !errors.Is(err, study.ErrInvalidCourse) {
		t.Fatalf("bad mode = %v", err)
	}
	w.add(study.ListTeachers, c, teacher.UserID)
	w.add(study.ListStudents, c, student.UserID)
	inherit, err := w.work.Create(w.ctx, teacher, c.ID, study.AssignmentInput{Title: "Парсер", Grading: study.GradingPass, Published: true})
	if err != nil {
		t.Fatalf("assignment: %v", err)
	}
	studyMode := module.ModeStudy
	own, err := w.work.Create(w.ctx, teacher, c.ID, study.AssignmentInput{Title: "Эссе", Grading: study.GradingPass, Published: true, NavigatorMode: &studyMode})
	if err != nil {
		t.Fatalf("assignment: %v", err)
	}
	briefer := study.Module{}.Briefers(module.Kernel{Tenants: w.tx, Access: w.access, Notifier: w.sent})[0]
	for _, tc := range []struct {
		id, mode string
	}{{inherit.ID, module.ModeProject}, {own.ID, module.ModeStudy}} {
		err := w.tx.Run(w.ctx, w.tenant, func(tx pgx.Tx) error {
			b, err := briefer.Brief(w.ctx, tx, student, tc.id)
			if err == nil && (b.Mode != tc.mode || b.Course != "Программирование" || b.CourseID != c.ID || b.Teaches) {
				t.Errorf("brief %s = %+v", tc.id, b)
			}
			return err
		})
		if err != nil {
			t.Fatalf("brief: %v", err)
		}
	}
	err = w.tx.Run(w.ctx, w.tenant, func(tx pgx.Tx) error {
		_, err := briefer.Brief(w.ctx, tx, outsider, own.ID)
		return err
	})
	if !errors.Is(err, module.ErrBriefNotFound) {
		t.Fatalf("outsider brief = %v", err)
	}
}

func TestCourseBriefTellsWhoTeaches(t *testing.T) {
	w := newWorld(t)
	teacher, student, outsider := w.user("dana", "dana@kaznu.kz"), w.user("aigerim", "aigerim@kaznu.kz"), w.user("erlan", "erlan@kaznu.kz")
	c, err := w.study.Create(w.ctx, w.admin, study.CourseInput{Title: "Биология"})
	if err != nil {
		t.Fatalf("course: %v", err)
	}
	w.add(study.ListTeachers, c, teacher.UserID)
	w.add(study.ListStudents, c, student.UserID)
	briefer := study.Module{}.Briefers(module.Kernel{Tenants: w.tx, Access: w.access, Notifier: w.sent})[1]
	if briefer.Kind() != study.BriefCourse {
		t.Fatalf("kind = %s", briefer.Kind())
	}
	err = w.tx.Run(w.ctx, w.tenant, func(tx pgx.Tx) error {
		taught, err := briefer.Brief(w.ctx, tx, teacher, c.ID)
		if err != nil {
			return err
		}
		learnt, err := briefer.Brief(w.ctx, tx, student, c.ID)
		if err != nil {
			return err
		}
		if !taught.Teaches || learnt.Teaches || learnt.CourseID != c.ID || learnt.Mode != module.ModeStudy {
			t.Errorf("teacher %+v, student %+v", taught, learnt)
		}
		_, err = briefer.Brief(w.ctx, tx, outsider, c.ID)
		if !errors.Is(err, module.ErrBriefNotFound) {
			t.Errorf("outsider = %v", err)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("brief: %v", err)
	}
}
