package structure_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/structure"
)

func TestTreeAndPath(t *testing.T) {
	w := newWorld(t)
	faculty := w.unit(nil, "faculty", "Физтех")
	dept := w.unit(&faculty, "department", "Кафедра ИИ")
	group := w.unit(&dept, "group", "ИИ-21")

	d, err := w.structure.Get(w.ctx, w.admin, group.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if len(d.Path) != 2 || d.Path[0].ID != faculty.ID || d.Path[1].ID != dept.ID {
		t.Fatalf("path = %+v", d.Path)
	}
	roots, err := w.structure.Children(w.ctx, w.admin, nil, false)
	if err != nil || len(roots) != 1 {
		t.Fatalf("roots = %+v, %v", roots, err)
	}
}

func TestRoleOnFacultyWorksInsideIt(t *testing.T) {
	w := newWorld(t)
	physics := w.unit(nil, "faculty", "Физтех")
	law := w.unit(nil, "faculty", "Юрфак")
	dean := w.user("dean", "dean@kaznu.kz")
	w.grantAt(dean, physics, structure.ManageStructure)

	if _, err := w.structure.Create(w.ctx, dean, structure.UnitInput{ParentID: &physics.ID, Kind: "group", Name: "Ф-1"}); err != nil {
		t.Fatalf("dean creates group in own faculty: %v", err)
	}
	if _, err := w.structure.Create(w.ctx, dean, structure.UnitInput{ParentID: &law.ID, Kind: "group", Name: "Ю-1"}); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("dean in other faculty: err = %v", err)
	}
	if _, err := w.structure.Create(w.ctx, dean, structure.UnitInput{Kind: "faculty", Name: "Новый"}); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("dean creates faculty: err = %v", err)
	}
}

func TestMoveRejectsCycle(t *testing.T) {
	w := newWorld(t)
	a := w.unit(nil, "faculty", "A")
	b := w.unit(&a, "department", "B")

	if err := w.structure.Move(w.ctx, w.admin, a.ID, &b.ID); !errors.Is(err, structure.ErrCycle) {
		t.Fatalf("err = %v, want ErrCycle", err)
	}
	if err := w.structure.Move(w.ctx, w.admin, b.ID, nil); err != nil {
		t.Fatalf("move to root: %v", err)
	}
}

func TestArchiveHidesSubtree(t *testing.T) {
	w := newWorld(t)
	a := w.unit(nil, "faculty", "A")
	b := w.unit(&a, "department", "B")

	if err := w.structure.SetArchived(w.ctx, w.admin, a.ID, true); err != nil {
		t.Fatalf("archive: %v", err)
	}
	children, err := w.structure.Children(w.ctx, w.admin, &a.ID, false)
	if err != nil || len(children) != 0 {
		t.Fatalf("children = %+v, %v", children, err)
	}
	if err := w.structure.SetArchived(w.ctx, w.admin, b.ID, false); !errors.Is(err, structure.ErrArchived) {
		t.Fatalf("restore under archived parent: err = %v", err)
	}
	if _, err := w.structure.Create(w.ctx, w.admin, structure.UnitInput{ParentID: &a.ID, Kind: "group", Name: "x"}); !errors.Is(err, structure.ErrArchived) {
		t.Fatalf("create under archived: err = %v", err)
	}
}

func TestMembersAndHistory(t *testing.T) {
	w := newWorld(t)
	group := w.unit(nil, "group", "ИИ-21")
	student := w.user("s1", "s1@kaznu.kz")

	if err := w.structure.AddMember(w.ctx, w.admin, group.ID, student.UserID); err != nil {
		t.Fatalf("add: %v", err)
	}
	if err := w.structure.AddMember(w.ctx, w.admin, group.ID, student.UserID); err != nil {
		t.Fatalf("add twice: %v", err)
	}
	members, err := w.structure.Members(w.ctx, w.admin, group.ID)
	if err != nil || len(members) != 1 || members[0].Email != "s1@kaznu.kz" {
		t.Fatalf("members = %+v, %v", members, err)
	}
	units, err := w.structure.UserUnits(w.ctx, student, student.UserID)
	if err != nil || len(units) != 1 {
		t.Fatalf("user units = %+v, %v", units, err)
	}
	if err := w.structure.RemoveMember(w.ctx, w.admin, group.ID, student.UserID); err != nil {
		t.Fatalf("remove: %v", err)
	}
	history, err := w.structure.History(w.ctx, w.admin, group.ID, 0)
	if err != nil {
		t.Fatalf("history: %v", err)
	}
	var actions []string
	for _, e := range history {
		actions = append(actions, e.Action)
	}
	want := []string{structure.ActionMemberRemoved, structure.ActionMemberAdded, structure.ActionCreated}
	if strings.Join(actions, ",") != strings.Join(want, ",") {
		t.Fatalf("actions = %v, want %v", actions, want)
	}
	if err := w.structure.AddMember(w.ctx, w.admin, group.ID, "not-a-user"); !errors.Is(err, structure.ErrUserNotFound) {
		t.Fatalf("unknown user: err = %v", err)
	}
}

func TestImportUnitsIsAllOrNothing(t *testing.T) {
	w := newWorld(t)
	bad := "code,parent_code,kind,name\nF1,,faculty,Физтех\nG1,NOPE,group,ИИ-21\nG2,F1,Group,Плохой\n"

	report, err := w.structure.ImportUnits(w.ctx, w.admin, strings.NewReader(bad))
	if err != nil {
		t.Fatalf("import: %v", err)
	}
	if len(report.Errors) != 2 || report.Errors[0].Line != 3 || report.Errors[0].Code != structure.ImportMissingParent || report.Errors[1].Code != structure.ImportInvalidUnit {
		t.Fatalf("report = %+v", report)
	}
	if roots, _ := w.structure.Children(w.ctx, w.admin, nil, true); len(roots) != 0 {
		t.Fatalf("failed import left units: %+v", roots)
	}

	good := "code,parent_code,kind,name\nF1,,faculty,Физтех\nG1,F1,group,ИИ-21\n"
	report, err = w.structure.ImportUnits(w.ctx, w.admin, strings.NewReader(good))
	if err != nil || report.Created != 2 || len(report.Errors) != 0 {
		t.Fatalf("report = %+v, %v", report, err)
	}
	again := "code,parent_code,kind,name\nG1,,group,ИИ-21 новая\n"
	report, err = w.structure.ImportUnits(w.ctx, w.admin, strings.NewReader(again))
	if err != nil || report.Updated != 1 {
		t.Fatalf("reimport = %+v, %v", report, err)
	}
	roots, _ := w.structure.Children(w.ctx, w.admin, nil, false)
	if len(roots) != 2 {
		t.Fatalf("G1 was not moved to root: %+v", roots)
	}
}

func TestImportMembers(t *testing.T) {
	w := newWorld(t)
	if _, err := w.structure.ImportUnits(w.ctx, w.admin, strings.NewReader("code,parent_code,kind,name\nG1,,group,ИИ-21\n")); err != nil {
		t.Fatalf("units: %v", err)
	}
	w.user("s1", "s1@kaznu.kz")

	report, err := w.structure.ImportMembers(w.ctx, w.admin, strings.NewReader("unit_code,email\nG1,S1@kaznu.kz\nG1,ghost@kaznu.kz\n"))
	if err != nil || len(report.Errors) != 1 || report.Errors[0].Code != structure.ImportUnknownUser {
		t.Fatalf("report = %+v, %v", report, err)
	}
	report, err = w.structure.ImportMembers(w.ctx, w.admin, strings.NewReader("unit_code,email\nG1,S1@kaznu.kz\n"))
	if err != nil || report.Created != 1 {
		t.Fatalf("report = %+v, %v", report, err)
	}
}

func TestImportRejectsMissingColumns(t *testing.T) {
	w := newWorld(t)

	if _, err := w.structure.ImportUnits(w.ctx, w.admin, strings.NewReader("code,name\nA,B\n")); !errors.Is(err, structure.ErrBadImport) {
		t.Fatalf("err = %v, want ErrBadImport", err)
	}
}

func TestUnknownUnitScopeIsInvalid(t *testing.T) {
	w := newWorld(t)

	if _, err := w.access.Permissions(w.ctx, w.admin, "unit:00000000-0000-0000-0000-000000000000"); !errors.Is(err, access.ErrInvalidScope) {
		t.Fatalf("err = %v, want ErrInvalidScope", err)
	}
}
