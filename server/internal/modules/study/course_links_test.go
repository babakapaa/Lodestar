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
		return nil
	})
	if err != nil {
		t.Fatalf("tx: %v", err)
	}
}
