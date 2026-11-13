package study_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/modules/study"
	"github.com/babakapaa/lodestar/server/internal/platform/table"
)

func TestGradebookAndStudentImport(t *testing.T) {
	w := newWorld(t)
	teacher := w.user("teacher", "teacher@kaznu.kz")
	aigerim := w.user("aigerim", "aigerim@kaznu.kz")
	daniyar := w.user("daniyar", "daniyar@kaznu.kz")
	c := w.course(nil, "Биология")
	w.add(study.ListTeachers, c, teacher.UserID)
	imports := study.NewStudentImportService(w.tx, w.access)
	books := study.NewGradebookService(w.tx, w.access)

	var file bytes.Buffer
	if err := table.Write(&file, table.XLSX, []string{"Email", "Name"}, [][]any{{"AIGERIM@kaznu.kz", "Айгерим"}, {"daniyar@kaznu.kz", ""}}); err != nil {
		t.Fatalf("write: %v", err)
	}
	report, err := imports.Import(w.ctx, teacher, c.ID, &file)
	if err != nil || report.Added != 2 || len(report.Errors) != 0 {
		t.Fatalf("import = %+v, %v", report, err)
	}
	report, err = imports.Import(w.ctx, teacher, c.ID, strings.NewReader("email,name\naigerim@kaznu.kz\nghost@kaznu.kz\n\n,Икс\n"))
	if err != nil || report.Added != 0 || len(report.Errors) != 2 || report.Errors[0].Line != 3 || report.Errors[0].Code != study.ImportUnknownUser || report.Errors[1].Code != study.ImportMissingEmail || report.Errors[1].Line != 5 {
		t.Fatalf("bad import = %+v, %v", report, err)
	}
	if _, err := imports.Import(w.ctx, teacher, c.ID, strings.NewReader("name\nx\n")); !errors.Is(err, study.ErrBadImport) {
		t.Fatalf("no email column err = %v", err)
	}
	if _, err := imports.Import(w.ctx, aigerim, c.ID, strings.NewReader("email\nx@kaznu.kz\n")); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("student import err = %v", err)
	}

	maxPoints := 20.0
	a, err := w.work.Create(w.ctx, teacher, c.ID, study.AssignmentInput{Title: "Эссе", Grading: study.GradingPoints, MaxPoints: &maxPoints, Published: true})
	if err != nil {
		t.Fatalf("assignment: %v", err)
	}
	if _, err := w.work.Create(w.ctx, teacher, c.ID, study.AssignmentInput{Title: "Чтение", Published: true}); err != nil {
		t.Fatalf("ungraded: %v", err)
	}
	points := 17.5
	if _, err := w.work.Grade(w.ctx, teacher, a.ID, aigerim.UserID, study.GradeInput{Points: &points}); err != nil {
		t.Fatalf("grade: %v", err)
	}
	book, err := books.Build(w.ctx, teacher, c.ID)
	if err != nil {
		t.Fatalf("gradebook: %v", err)
	}
	if strings.Join(book.Header, "|") != "name|email|Эссе (20)" || len(book.Rows) != 2 {
		t.Fatalf("book = %+v", book)
	}
	got := map[string]any{}
	for _, row := range book.Rows {
		got[row[1].(string)] = row[2]
	}
	if got["aigerim@kaznu.kz"] != 17.5 || got["daniyar@kaznu.kz"] != nil {
		t.Fatalf("scores = %+v", got)
	}
	if _, err := books.Build(w.ctx, daniyar, c.ID); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("student gradebook err = %v", err)
	}
}
