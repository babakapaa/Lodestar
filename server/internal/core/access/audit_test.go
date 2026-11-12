package access_test

import (
	"errors"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/audit"
)

func TestRoleChangesAreLogged(t *testing.T) {
	w := newWorld(t)
	role := w.role("Куратор", postsWrite)
	student := w.user("student")
	if _, err := w.service.Assign(w.ctx, w.admin, access.AssignInput{RoleID: role.ID, UserID: student.UserID}); err != nil {
		t.Fatalf("assign: %v", err)
	}
	log := audit.NewService(w.tx, w.service)

	entries, err := log.List(w.ctx, w.admin, audit.Filter{Action: "core.role"})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(entries) != 2 || entries[0].Action != access.ActionRoleAssigned || entries[1].Action != access.ActionRoleCreated {
		t.Fatalf("entries = %+v", entries)
	}
	if entries[0].ActorID == nil || *entries[0].ActorID != w.admin.UserID || entries[0].TargetID != student.UserID {
		t.Fatalf("assigned entry = %+v", entries[0])
	}
	older, _ := log.List(w.ctx, w.admin, audit.Filter{Before: entries[0].ID})
	if len(older) != 1 || older[0].ID != entries[1].ID {
		t.Fatalf("older = %+v", older)
	}
}

func TestAuditLogNeedsPermission(t *testing.T) {
	w := newWorld(t)
	if _, err := audit.NewService(w.tx, w.service).List(w.ctx, w.user("curious"), audit.Filter{}); !errors.Is(err, audit.ErrForbidden) {
		t.Fatalf("err = %v, want ErrForbidden", err)
	}
}
