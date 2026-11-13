package study_test

import (
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/files"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/modules/study"
)

func TestAssignmentFromPublicationToGrade(t *testing.T) {
	w := newWorld(t)
	group := w.unit(nil, "group", "ИС-21")
	teacher := w.user("dana", "dana@kaznu.kz")
	student := w.user("aigerim", "aigerim@kaznu.kz")
	outsider := w.user("erlan", "erlan@kaznu.kz")
	w.join(group, student)
	c := w.course(nil, "Базы данных")
	w.add(study.ListUnits, c, group.ID)
	w.add(study.ListTeachers, c, teacher.UserID)

	brief := w.file(teacher, "brief.txt")
	due := time.Now().Add(48 * time.Hour)
	maxPoints := 10.0
	draft, err := w.work.Create(w.ctx, teacher, c.ID, study.AssignmentInput{
		Title: "Нормальные формы", DueAt: &due, Grading: study.GradingPoints, MaxPoints: &maxPoints, Files: []string{brief},
	})
	if err != nil {
		t.Fatalf("create draft: %v", err)
	}
	if list, _ := w.work.List(w.ctx, student, c.ID); len(list) != 0 {
		t.Fatalf("student sees draft: %+v", list)
	}
	if _, err := w.files.Get(w.ctx, student, brief); !errors.Is(err, files.ErrNotFound) {
		t.Fatalf("draft file readable: err = %v", err)
	}
	if len(w.sent.of(study.NotifyPublished)) != 0 {
		t.Fatal("draft announced")
	}

	published, err := w.work.Update(w.ctx, teacher, draft.ID, study.AssignmentInput{
		Title: "Нормальные формы", DueAt: &due, Grading: study.GradingPoints, MaxPoints: &maxPoints, Files: []string{brief}, Published: true,
	})
	if err != nil || !published.Published() {
		t.Fatalf("publish: %+v, %v", published, err)
	}
	announced := w.sent.of(study.NotifyPublished)
	if len(announced) != 1 || len(announced[0].UserIDs) != 1 || announced[0].UserIDs[0] != student.UserID {
		t.Fatalf("announcement = %+v", announced)
	}
	if _, err := w.files.Get(w.ctx, student, brief); err != nil {
		t.Fatalf("student reads brief: %v", err)
	}
	if _, err := w.files.Get(w.ctx, outsider, brief); !errors.Is(err, files.ErrNotFound) {
		t.Fatalf("outsider reads brief: err = %v", err)
	}

	answer := w.file(student, "answer.txt")
	composed := time.Now().Add(-3 * time.Hour)
	in := study.SubmissionInput{ClientID: "6f1d7a52-3c38-4a8e-9a43-1f0e5f3c1a01", Body: "Ответ", Files: []string{answer}, ComposedAt: &composed}
	sub, created, err := w.work.Submit(w.ctx, student, published.ID, in)
	if err != nil || !created || len(sub.Files) != 1 {
		t.Fatalf("submit: %+v, %v, %v", sub, created, err)
	}
	again, created, err := w.work.Submit(w.ctx, student, published.ID, in)
	if err != nil || created || again.ID != sub.ID {
		t.Fatalf("resend: %+v, %v, %v", again, created, err)
	}
	if n := len(w.sent.of(study.NotifySubmitted)); n != 1 {
		t.Fatalf("teacher told %d times", n)
	}
	if _, err := w.files.Get(w.ctx, teacher, answer); err != nil {
		t.Fatalf("teacher reads answer: %v", err)
	}
	if _, err := w.files.Get(w.ctx, outsider, answer); !errors.Is(err, files.ErrNotFound) {
		t.Fatalf("outsider reads answer: err = %v", err)
	}

	work, err := w.work.Work(w.ctx, teacher, published.ID)
	if err != nil || len(work) != 1 || work[0].Latest == nil || work[0].Attempts != 1 {
		t.Fatalf("work = %+v, %v", work, err)
	}
	tooMuch := 11.0
	if _, err := w.work.Grade(w.ctx, teacher, published.ID, student.UserID, study.GradeInput{Points: &tooMuch}); !errors.Is(err, study.ErrInvalidGrade) {
		t.Fatalf("too many points: err = %v", err)
	}
	if _, err := w.work.Grade(w.ctx, teacher, published.ID, outsider.UserID, study.GradeInput{Points: &maxPoints}); !errors.Is(err, study.ErrNotStudent) {
		t.Fatalf("grade outsider: err = %v", err)
	}
	if _, err := w.work.Grade(w.ctx, student, published.ID, student.UserID, study.GradeInput{Points: &maxPoints}); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("student grades self: err = %v", err)
	}
	nine := 9.0
	if _, err := w.work.Grade(w.ctx, teacher, published.ID, student.UserID, study.GradeInput{Points: &nine, Comment: "Хорошо"}); err != nil {
		t.Fatalf("grade: %v", err)
	}
	graded := w.sent.of(study.NotifyGraded)
	if len(graded) != 1 || graded[0].UserIDs[0] != student.UserID || graded[0].Link != "/study/assignments/"+published.ID {
		t.Fatalf("graded = %+v", graded)
	}

	list, err := w.work.List(w.ctx, student, c.ID)
	if err != nil || len(list) != 1 || list[0].Mine == nil || list[0].Mine.SubmittedAt == nil || list[0].Mine.Grade == nil || *list[0].Mine.Grade.Points != 9 {
		t.Fatalf("student list = %+v, %v", list, err)
	}
	details, err := w.work.Get(w.ctx, student, published.ID)
	if err != nil || details.Mine == nil || details.Mine.Grade == nil || len(details.Files) != 1 || details.CourseTitle != "Базы данных" {
		t.Fatalf("student details = %+v, %v", details, err)
	}
	teaching, err := w.work.List(w.ctx, teacher, c.ID)
	if err != nil || teaching[0].Counts == nil || *teaching[0].Counts != (study.Counts{Students: 1, Submitted: 1, Graded: 1}) {
		t.Fatalf("teacher list = %+v, %v", teaching, err)
	}
	if err := w.work.Delete(w.ctx, teacher, published.ID); !errors.Is(err, study.ErrHasWork) {
		t.Fatalf("delete with work: err = %v", err)
	}
	if _, err := w.work.Attempts(w.ctx, outsider, published.ID, student.UserID); !errors.Is(err, study.ErrAssignmentNotFound) {
		t.Fatalf("outsider attempts: err = %v", err)
	}
	var metrics []module.Metric
	err = w.tx.Run(w.ctx, w.tenant, func(tx pgx.Tx) error {
		var err error
		metrics, err = study.StudyStats{}.Stats(w.ctx, tx, []string{student.UserID, outsider.UserID}, time.Now().Add(-time.Hour), time.Now().Add(time.Hour))
		return err
	})
	want := []module.Metric{{Key: study.MetricSubmitted, Value: 0.5}, {Key: study.MetricOnTime, Value: 1}}
	if err != nil || !slices.Equal(metrics, want) {
		t.Fatalf("stats = %+v, %v", metrics, err)
	}
}

func TestSubmissionRules(t *testing.T) {
	w := newWorld(t)
	teacher := w.user("dana", "dana@kaznu.kz")
	student := w.user("aigerim", "aigerim@kaznu.kz")
	c := w.course(nil, "Физика")
	w.add(study.ListTeachers, c, teacher.UserID)
	w.add(study.ListStudents, c, student.UserID)
	a, err := w.work.Create(w.ctx, teacher, c.ID, study.AssignmentInput{Title: "Лаба", Grading: study.GradingPass, Published: true})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, _, err := w.work.Submit(w.ctx, student, a.ID, study.SubmissionInput{ClientID: "6f1d7a52-3c38-4a8e-9a43-1f0e5f3c1a02"}); !errors.Is(err, study.ErrInvalidSubmission) {
		t.Fatalf("empty: err = %v", err)
	}
	foreign := w.file(teacher, "x.txt")
	if _, _, err := w.work.Submit(w.ctx, student, a.ID, study.SubmissionInput{ClientID: "6f1d7a52-3c38-4a8e-9a43-1f0e5f3c1a03", Files: []string{foreign}}); !errors.Is(err, study.ErrBadFile) {
		t.Fatalf("foreign file: err = %v", err)
	}
	if _, _, err := w.work.Submit(w.ctx, teacher, a.ID, study.SubmissionInput{ClientID: "6f1d7a52-3c38-4a8e-9a43-1f0e5f3c1a04", Body: "x"}); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("teacher submits: err = %v", err)
	}
	if _, _, err := w.work.Submit(w.ctx, student, a.ID, study.SubmissionInput{ClientID: "not-a-uuid", Body: "x"}); !errors.Is(err, study.ErrInvalidSubmission) {
		t.Fatalf("bad client id: err = %v", err)
	}
	passed := true
	points := 5.0
	if _, err := w.work.Grade(w.ctx, teacher, a.ID, student.UserID, study.GradeInput{Points: &points}); !errors.Is(err, study.ErrInvalidGrade) {
		t.Fatalf("points on pass grading: err = %v", err)
	}
	if _, err := w.work.Grade(w.ctx, teacher, a.ID, student.UserID, study.GradeInput{Passed: &passed}); err != nil {
		t.Fatalf("pass: %v", err)
	}
	if _, err := w.work.Create(w.ctx, teacher, c.ID, study.AssignmentInput{Title: "x", Grading: study.GradingPoints}); !errors.Is(err, study.ErrInvalidAssignment) {
		t.Fatalf("points without max: err = %v", err)
	}
	if _, err := w.work.Create(w.ctx, student, c.ID, study.AssignmentInput{Title: "x"}); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("student creates: err = %v", err)
	}
}
