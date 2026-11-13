package study_test

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/modules/study"
)

func TestCourseLinks(t *testing.T) {
	w := newWorld(t)
	teacher := w.user("teacher", "teacher@kaznu.kz")
	student := w.user("student", "student@kaznu.kz")
	bio := w.course(nil, "Биология")
	w.add(study.ListTeachers, bio, teacher.UserID)
	w.add(study.ListStudents, bio, student.UserID)
	source := study.Module{}.LinkSources(w.kernel())[0]
	if source.Kind() != study.LinkCourse {
		t.Fatalf("kind = %s", source.Kind())
	}
	err := w.tx.Run(w.ctx, w.tenant, func(tx pgx.Tx) error {
		list, err := source.Linkable(w.ctx, tx, teacher)
		if err != nil {
			return err
		}
		if len(list) != 1 || list[0].ID != bio.ID || list[0].Path != "/study/courses/"+bio.ID {
			t.Errorf("linkable = %+v", list)
		}
		if list, _ := source.Linkable(w.ctx, tx, student); len(list) != 0 {
			t.Errorf("student linkable = %+v", list)
		}
		if err := source.CanLink(w.ctx, tx, teacher, bio.ID); err != nil {
			t.Errorf("teacher can link: %v", err)
		}
		if err := source.CanLink(w.ctx, tx, student, bio.ID); !errors.Is(err, access.ErrForbidden) {
			t.Errorf("student can link: %v", err)
		}
		if _, err := source.Describe(w.ctx, tx, "00000000-0000-0000-0000-000000000000"); !errors.Is(err, module.ErrLinkTargetNotFound) {
			t.Errorf("describe missing: %v", err)
		}
		target, err := source.Describe(w.ctx, tx, bio.ID)
		if err != nil || target.Title != "Биология" {
			t.Errorf("describe = %+v, %v", target, err)
		}
		lead, learner := w.user("lead", "lead@kaznu.kz"), w.user("learner", "learner@kaznu.kz")
		added, err := source.Enroll(w.ctx, tx, teacher.UserID, bio.ID, []string{lead.UserID}, []string{learner.UserID, student.UserID})
		if err != nil || added != 2 {
			t.Errorf("enroll = %d, %v", added, err)
		}
		if added, _ := source.Enroll(w.ctx, tx, teacher.UserID, bio.ID, []string{lead.UserID}, []string{learner.UserID}); added != 0 {
			t.Errorf("second enroll = %d", added)
		}
		ref, err := source.PutWork(w.ctx, tx, teacher.UserID, bio.ID, "", module.ExternalWork{Title: "Эссе", Max: 20})
		if err != nil || ref.Kind != study.TeamScopeAssignment {
			t.Fatalf("put work = %+v, %v", ref, err)
		}
		again, err := source.PutWork(w.ctx, tx, teacher.UserID, bio.ID, ref.ID, module.ExternalWork{Title: "Эссе 2", Max: 25})
		if err != nil || again.ID != ref.ID {
			t.Errorf("update work = %+v, %v", again, err)
		}
		var title string
		var max float64
		if err := tx.QueryRow(w.ctx, `SELECT title, max_points::float8 FROM study.assignments WHERE id::text = $1`, ref.ID).Scan(&title, &max); err != nil || title != "Эссе 2" || max != 25 {
			t.Errorf("assignment = %s %v %v", title, max, err)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("tx: %v", err)
	}
}
