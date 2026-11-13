package study_test

import (
	"errors"
	"testing"
	"time"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/modules/study"
)

func TestTeacherSeesCourseThroughStudentEyes(t *testing.T) {
	w := newWorld(t)
	group := w.unit(nil, "group", "ИС-21")
	teacher := w.user("dana", "dana@kaznu.kz")
	student := w.user("aigerim", "aigerim@kaznu.kz")
	outsider := w.user("erlan", "erlan@kaznu.kz")
	w.join(group, student)
	c := w.course(nil, "Базы данных")
	w.add(study.ListUnits, c, group.ID)
	w.add(study.ListTeachers, c, teacher.UserID)

	due := time.Now().Add(48 * time.Hour)
	if _, err := w.work.Create(w.ctx, teacher, c.ID, study.AssignmentInput{Title: "Черновик", DueAt: &due, Grading: study.GradingPass}); err != nil {
		t.Fatalf("draft: %v", err)
	}
	open, err := w.work.Create(w.ctx, teacher, c.ID, study.AssignmentInput{Title: "Нормальные формы", DueAt: &due, Grading: study.GradingPass, Published: true})
	if err != nil {
		t.Fatalf("publish: %v", err)
	}
	if _, _, err := w.work.Submit(w.ctx, student, open.ID, study.SubmissionInput{ClientID: "6f1d7a52-3c38-4a8e-9a43-1f0e5f3c1a02", Body: "Ответ"}); err != nil {
		t.Fatalf("submit: %v", err)
	}

	later := time.Now().Add(time.Hour)
	seen, err := w.work.Preview(w.ctx, teacher, c.ID, student.UserID, later)
	if err != nil || len(seen) != 1 || seen[0].ID != open.ID || seen[0].Mine == nil || seen[0].Mine.SubmittedAt == nil || seen[0].Counts != nil {
		t.Fatalf("preview = %+v, %v", seen, err)
	}
	anyone, err := w.work.Preview(w.ctx, teacher, c.ID, "", later)
	if err != nil || len(anyone) != 1 || anyone[0].Mine.SubmittedAt != nil {
		t.Fatalf("preview without student = %+v, %v", anyone, err)
	}
	if past, err := w.work.Preview(w.ctx, teacher, c.ID, student.UserID, time.Now().Add(-time.Hour)); err != nil || len(past) != 0 {
		t.Fatalf("preview in the past = %+v, %v", past, err)
	}
	if _, err := w.work.Preview(w.ctx, teacher, c.ID, outsider.UserID, later); !errors.Is(err, study.ErrNotStudent) {
		t.Fatalf("outsider: %v", err)
	}
	if _, err := w.work.Preview(w.ctx, student, c.ID, student.UserID, later); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("student previews: %v", err)
	}
}
