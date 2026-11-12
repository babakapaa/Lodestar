package access_test

import (
	"errors"
	"testing"
	"time"

	"github.com/babakapaa/lodestar/server/internal/core/access"
)

func TestBootstrapGivesTopRoleEverything(t *testing.T) {
	w := newWorld(t)

	if !w.can(w.admin, "", postsPin) || !w.can(w.admin, "", access.ManageRoles) {
		t.Fatal("administrator is missing permissions")
	}
	roles, err := w.service.Roles(w.ctx, w.admin)
	if err != nil {
		t.Fatalf("roles: %v", err)
	}
	if len(roles) != 2 || roles[0].Name != access.AdminRoleName || !roles[1].Default {
		t.Fatalf("roles = %+v", roles)
	}
}

func TestRolesAddUp(t *testing.T) {
	w := newWorld(t)
	student := w.user("student")
	if w.can(student, "", postsWrite) {
		t.Fatal("new user already has a permission")
	}
	w.assign(student, w.role("Автор", postsWrite), "")
	w.assign(student, w.role("Редактор", postsPin), "")

	if !w.can(student, "", postsWrite) || !w.can(student, "", postsPin) {
		t.Fatal("permissions of two roles did not add up")
	}
}

func TestDefaultRoleAppliesToEveryone(t *testing.T) {
	w := newWorld(t)
	def := w.defaultRole()
	if _, err := w.service.UpdateRole(w.ctx, w.admin, def.ID, access.RoleInput{Name: def.Name, Permissions: []access.Permission{postsWrite}}); err != nil {
		t.Fatalf("update default: %v", err)
	}

	if !w.can(w.user("anyone"), "", postsWrite) {
		t.Fatal("default role permission is not applied")
	}
}

func TestScopedRoleWorksOnlyInsideScope(t *testing.T) {
	w := newWorld(t)
	head := w.user("head")
	w.assign(head, w.role("Староста", postsPin), "group:g1")

	if !w.can(head, "group:g1", postsPin) {
		t.Fatal("scoped role does not work in its group")
	}
	if w.can(head, "group:g2", postsPin) || w.can(head, "", postsPin) {
		t.Fatal("scoped role leaks outside its group")
	}
}

func TestDeeperOverrideWins(t *testing.T) {
	w := newWorld(t)
	author := w.user("author")
	role := w.role("Автор", postsWrite)
	w.assign(author, role, "")

	deny := access.Override{RoleID: role.ID, Deny: []access.Permission{postsWrite}}
	if err := w.service.SetOverride(w.ctx, w.admin, "faculty:f1", deny); err != nil {
		t.Fatalf("deny: %v", err)
	}
	allow := access.Override{RoleID: role.ID, Allow: []access.Permission{postsWrite}}
	if err := w.service.SetOverride(w.ctx, w.admin, "group:g1", allow); err != nil {
		t.Fatalf("allow: %v", err)
	}

	if w.can(author, "faculty:f1", postsWrite) {
		t.Fatal("deny on faculty is ignored")
	}
	if w.can(author, "group:g2", postsWrite) {
		t.Fatal("deny on faculty is not inherited by its groups")
	}
	if !w.can(author, "group:g1", postsWrite) {
		t.Fatal("allow on group does not beat deny on faculty")
	}
}

func TestOverrideOnTenantIsRejected(t *testing.T) {
	w := newWorld(t)
	role := w.role("Автор")

	err := w.service.SetOverride(w.ctx, w.admin, "", access.Override{RoleID: role.ID, Deny: []access.Permission{postsWrite}})
	if !errors.Is(err, access.ErrInvalidScope) {
		t.Fatalf("err = %v, want ErrInvalidScope", err)
	}
}

func TestExpiredRoleStopsWorking(t *testing.T) {
	w := newWorld(t)
	guest := w.user("guest")
	past := time.Now().Add(-time.Hour)
	_, err := w.service.Assign(w.ctx, w.admin, access.AssignInput{RoleID: w.role("Гость", postsWrite).ID, UserID: guest.UserID, ExpiresAt: &past})
	if err != nil {
		t.Fatalf("assign: %v", err)
	}

	if w.can(guest, "", postsWrite) {
		t.Fatal("expired role still grants permissions")
	}
}

func TestManagerIsBoundByHierarchy(t *testing.T) {
	w := newWorld(t)
	moderator := w.user("moderator")
	modRole := w.role("Модератор", access.ManageRoles, postsWrite)
	w.assign(moderator, modRole, "")
	student := w.user("student")

	below, err := w.service.CreateRole(w.ctx, moderator, access.RoleInput{Name: "Автор", Permissions: []access.Permission{postsWrite}})
	if err != nil {
		t.Fatalf("moderator creates role: %v", err)
	}
	if below.Position >= modRole.Position+1 {
		t.Fatalf("new role position %d is not below moderator", below.Position)
	}
	if _, err := w.service.CreateRole(w.ctx, moderator, access.RoleInput{Name: "x", Permissions: []access.Permission{postsPin}}); !errors.Is(err, access.ErrPermissionNotHeld) {
		t.Fatalf("grant not held: err = %v", err)
	}
	if _, err := w.service.UpdateRole(w.ctx, moderator, modRole.ID, access.RoleInput{Name: "Главный"}); !errors.Is(err, access.ErrHierarchy) {
		t.Fatalf("edit own role: err = %v", err)
	}
	if _, err := w.service.Assign(w.ctx, moderator, access.AssignInput{RoleID: modRole.ID, UserID: student.UserID}); !errors.Is(err, access.ErrHierarchy) {
		t.Fatalf("assign equal role: err = %v", err)
	}
	if _, err := w.service.Assign(w.ctx, moderator, access.AssignInput{RoleID: below.ID, UserID: student.UserID}); err != nil {
		t.Fatalf("assign lower role: %v", err)
	}
	if _, err := w.service.CreateRole(w.ctx, student, access.RoleInput{Name: "y"}); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("student creates role: err = %v", err)
	}
}

func TestMoveAndDeleteKeepOrder(t *testing.T) {
	w := newWorld(t)
	a := w.role("A")
	w.role("B")
	c := w.role("C")

	if err := w.service.MoveRole(w.ctx, w.admin, a.ID, c.Position); err != nil {
		t.Fatalf("move: %v", err)
	}
	if err := w.service.DeleteRole(w.ctx, w.admin, c.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	roles, err := w.service.Roles(w.ctx, w.admin)
	if err != nil {
		t.Fatalf("roles: %v", err)
	}
	var names []string
	for i, r := range roles {
		names = append(names, r.Name)
		if r.Position != len(roles)-1-i {
			t.Fatalf("positions have gaps: %+v", roles)
		}
	}
	want := []string{access.AdminRoleName, "A", "B", access.DefaultRoleName}
	if len(names) != len(want) {
		t.Fatalf("names = %v, want %v", names, want)
	}
	for i := range want {
		if names[i] != want[i] {
			t.Fatalf("names = %v, want %v", names, want)
		}
	}
	if err := w.service.DeleteRole(w.ctx, w.admin, w.defaultRole().ID); !errors.Is(err, access.ErrHierarchy) && !errors.Is(err, access.ErrDefaultRole) {
		t.Fatalf("delete default: err = %v", err)
	}
}

func TestPreviewShowsRoleAlone(t *testing.T) {
	w := newWorld(t)
	role := w.role("Староста", postsPin)

	perms, err := w.service.Preview(w.ctx, w.admin, role.ID, "group:g1")
	if err != nil {
		t.Fatalf("preview: %v", err)
	}
	if len(perms) != 1 || perms[0] != postsPin {
		t.Fatalf("preview = %v", perms)
	}
}

func TestUnknownPermissionIsRejected(t *testing.T) {
	w := newWorld(t)

	_, err := w.service.CreateRole(w.ctx, w.admin, access.RoleInput{Name: "x", Permissions: []access.Permission{"made.up"}})
	if !errors.Is(err, access.ErrUnknownPermission) {
		t.Fatalf("err = %v, want ErrUnknownPermission", err)
	}
}

func TestUnassignRemovesRole(t *testing.T) {
	w := newWorld(t)
	student := w.user("student")
	a := w.assign(student, w.role("Автор", postsWrite), "")

	if err := w.service.Unassign(w.ctx, w.admin, a.ID); err != nil {
		t.Fatalf("unassign: %v", err)
	}
	if w.can(student, "", postsWrite) {
		t.Fatal("permission survived unassign")
	}
}
