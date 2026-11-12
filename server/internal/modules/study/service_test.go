package study_test

import (
	"errors"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/modules/study"
)

func TestStudentsComeFromGroupsInsideTheStream(t *testing.T) {
	w := newWorld(t)
	stream := w.unit(nil, "stream", "ИС-2023")
	group := w.unit(&stream, "group", "ИС-21")
	other := w.unit(nil, "group", "ЮР-21")
	student := w.user("aigerim", "aigerim@kaznu.kz")
	outsider := w.user("erlan", "erlan@kaznu.kz")
	teacher := w.user("dana", "dana@kaznu.kz")
	w.join(group, student)
	w.join(other, outsider)
	c := w.course(nil, "Базы данных")
	w.add(study.ListUnits, c, stream.ID)
	w.add(study.ListTeachers, c, teacher.UserID)

	mine, err := w.study.Mine(w.ctx, student, false)
	if err != nil || len(mine) != 1 || mine[0].Role != study.RoleStudent {
		t.Fatalf("student courses = %+v, %v", mine, err)
	}
	taught, err := w.study.Mine(w.ctx, teacher, false)
	if err != nil || len(taught) != 1 || taught[0].Role != study.RoleTeacher {
		t.Fatalf("teacher courses = %+v, %v", taught, err)
	}
	if none, _ := w.study.Mine(w.ctx, outsider, false); len(none) != 0 {
		t.Fatalf("outsider sees %+v", none)
	}
	if _, err := w.study.Get(w.ctx, outsider, c.ID); !errors.Is(err, study.ErrCourseNotFound) {
		t.Fatalf("outsider get: err = %v", err)
	}
	students, err := w.study.Students(w.ctx, teacher, c.ID)
	if err != nil || len(students) != 1 || students[0].UserID != student.UserID {
		t.Fatalf("students = %+v, %v", students, err)
	}
	d, err := w.study.Get(w.ctx, student, c.ID)
	if err != nil || d.Role != study.RoleStudent || d.CanManage || len(d.Teachers) != 1 || len(d.Units) != 1 {
		t.Fatalf("details = %+v, %v", d, err)
	}
}

func TestTeacherEditsCourseButNotItsOwner(t *testing.T) {
	w := newWorld(t)
	dept := w.unit(nil, "department", "Кафедра ИС")
	teacher := w.user("dana", "dana@kaznu.kz")
	student := w.user("aigerim", "aigerim@kaznu.kz")
	c := w.course(&dept, "Алгоритмы")
	w.add(study.ListTeachers, c, teacher.UserID)

	updated, err := w.study.Update(w.ctx, teacher, c.ID, study.CourseInput{UnitID: &dept.ID, Title: "Алгоритмы и структуры", Description: "Весна"})
	if err != nil || updated.Title != "Алгоритмы и структуры" {
		t.Fatalf("teacher update = %+v, %v", updated, err)
	}
	if _, err := w.study.Update(w.ctx, teacher, c.ID, study.CourseInput{Title: "Алгоритмы"}); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("teacher moves course: err = %v", err)
	}
	if err := w.study.Change(w.ctx, teacher, study.ListStudents, c.ID, student.UserID, true); err != nil {
		t.Fatalf("teacher adds student: %v", err)
	}
	if err := w.study.Change(w.ctx, teacher, study.ListTeachers, c.ID, student.UserID, true); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("teacher adds teacher: err = %v", err)
	}
	if err := w.study.SetArchived(w.ctx, teacher, c.ID, true); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("teacher archives: err = %v", err)
	}
	if _, err := w.study.Update(w.ctx, student, c.ID, study.CourseInput{UnitID: &dept.ID, Title: "x"}); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("student update: err = %v", err)
	}
}

func TestManagerOfDepartmentOwnsItsCourses(t *testing.T) {
	w := newWorld(t)
	dept := w.unit(nil, "department", "Кафедра ИС")
	law := w.unit(nil, "department", "Кафедра права")
	head := w.user("head", "head@kaznu.kz")
	w.grantAt(head, dept, study.ManageCourses)

	c, err := w.study.Create(w.ctx, head, study.CourseInput{UnitID: &dept.ID, Title: "Сети", Code: ptr("NET-1")})
	if err != nil {
		t.Fatalf("head creates course: %v", err)
	}
	if _, err := w.study.Create(w.ctx, head, study.CourseInput{UnitID: &law.ID, Title: "Право"}); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("head in other department: err = %v", err)
	}
	if _, err := w.study.Create(w.ctx, w.admin, study.CourseInput{Title: "Сети 2", Code: ptr("NET-1")}); !errors.Is(err, study.ErrCodeTaken) {
		t.Fatalf("duplicate code: err = %v", err)
	}
	list, err := w.study.Catalog(w.ctx, head, &dept.ID, false)
	if err != nil || len(list) != 1 {
		t.Fatalf("catalog = %+v, %v", list, err)
	}
	d, err := w.study.Get(w.ctx, head, c.ID)
	if err != nil || d.Role != study.RoleManager || !d.CanManage {
		t.Fatalf("head details = %+v, %v", d, err)
	}
	if err := w.study.SetArchived(w.ctx, head, c.ID, true); err != nil {
		t.Fatalf("archive: %v", err)
	}
	if err := w.study.Change(w.ctx, head, study.ListStudents, c.ID, head.UserID, true); !errors.Is(err, study.ErrArchived) {
		t.Fatalf("add to archived: err = %v", err)
	}
	if active, _ := w.study.Catalog(w.ctx, head, &dept.ID, false); len(active) != 0 {
		t.Fatalf("archived course listed: %+v", active)
	}
}

func TestInvalidInputIsRejected(t *testing.T) {
	w := newWorld(t)
	if _, err := w.study.Create(w.ctx, w.admin, study.CourseInput{Title: "  "}); !errors.Is(err, study.ErrInvalidCourse) {
		t.Fatalf("empty title: err = %v", err)
	}
	missing := "00000000-0000-0000-0000-000000000000"
	if _, err := w.study.Create(w.ctx, w.admin, study.CourseInput{UnitID: &missing, Title: "x"}); !errors.Is(err, study.ErrUnitNotFound) {
		t.Fatalf("missing unit: err = %v", err)
	}
	c := w.course(nil, "Физика")
	if err := w.study.Change(w.ctx, w.admin, study.ListTeachers, c.ID, "nope", true); !errors.Is(err, study.ErrUserNotFound) {
		t.Fatalf("bad user: err = %v", err)
	}
	if _, err := w.study.Get(w.ctx, w.admin, "nope"); !errors.Is(err, study.ErrCourseNotFound) {
		t.Fatalf("bad id: err = %v", err)
	}
}

func ptr(s string) *string {
	return &s
}
