package demo_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"slices"
	"testing"
	"time"

	"github.com/babakapaa/lodestar/server/internal/app"
	"github.com/babakapaa/lodestar/server/internal/app/demo"
	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/core/structure"
	"github.com/babakapaa/lodestar/server/internal/modules/chat"
	"github.com/babakapaa/lodestar/server/internal/modules/study"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/database/dbtest"
	"github.com/babakapaa/lodestar/server/internal/platform/secretbox"
)

func TestSeedIsRepeatableAndGivesEachPersonTheirRole(t *testing.T) {
	ctx := context.Background()
	pool := dbtest.New(t)
	registry, err := module.NewRegistry(app.Modules()...)
	if err != nil {
		t.Fatalf("registry: %v", err)
	}
	if err := database.NewMigrator(pool).Apply(ctx, append(app.CoreMigrations(), registry.Migrations()...)...); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	box, err := secretbox.New(base64.StdEncoding.EncodeToString(bytes.Repeat([]byte{5}, secretbox.KeySize)))
	if err != nil {
		t.Fatalf("secretbox: %v", err)
	}
	cfg := demo.Config{Issuer: "http://idp.test", ClientID: "lodestar-demo", ClientSecret: "secret", PublicURL: "http://lodestar.test"}

	first, err := demo.Seed(ctx, pool, box, registry, cfg)
	if err != nil {
		t.Fatalf("seed: %v", err)
	}
	second, err := demo.Seed(ctx, pool, box, registry, cfg)
	if err != nil {
		t.Fatalf("seed again: %v", err)
	}
	if len(second.People) != len(demo.People()) || second.People[demo.Peer.Subject] == "" {
		t.Fatalf("people = %v", second.People)
	}
	if first.Tenant.ID != second.Tenant.ID || second.Callback != "http://lodestar.test/v1/auth/demo/demo/callback" {
		t.Fatalf("results = %+v / %+v", first, second)
	}

	tx := database.NewTenantTx(pool)
	catalog, err := access.NewCatalog(app.CorePermissions()...)
	if err != nil {
		t.Fatalf("catalog: %v", err)
	}
	service := access.NewService(tx, catalog, structure.Tree{})
	principalOf := func(email string) auth.Principal {
		users, err := identity.NewStore(tx).ByEmail(ctx, first.Tenant.ID, email)
		if err != nil || len(users) != 1 {
			t.Fatalf("user %s: %v %d", email, err, len(users))
		}
		return auth.Principal{TenantID: first.Tenant.ID, UserID: users[0].ID}
	}
	permissionsOf := func(email string) []access.Permission {
		perms, err := service.Permissions(ctx, principalOf(email), "")
		if err != nil {
			t.Fatalf("permissions %s: %v", email, err)
		}
		return perms
	}
	if !slices.Contains(permissionsOf(demo.Admin.Email), access.Administrator) {
		t.Errorf("admin is not an administrator")
	}
	teacher := permissionsOf(demo.Teacher.Email)
	if !slices.Contains(teacher, structure.ManageMembers) || slices.Contains(teacher, access.Administrator) {
		t.Errorf("teacher permissions = %v", teacher)
	}
	if student := permissionsOf(demo.Student.Email); slices.Contains(student, structure.ManageMembers) {
		t.Errorf("student permissions = %v", student)
	}

	core, err := app.NewCore(pool, registry, box)
	if err != nil {
		t.Fatalf("core: %v", err)
	}
	student := principalOf(demo.Student.Email)
	courses, err := study.NewService(tx, core.Access).Mine(ctx, student, false)
	if err != nil || len(courses) != 1 || courses[0].Role != study.RoleStudent {
		t.Fatalf("student courses = %+v, %v", courses, err)
	}
	peerCourses, err := study.NewService(tx, core.Access).Mine(ctx, principalOf(demo.Peer.Email), false)
	if err != nil || len(peerCourses) != 1 || peerCourses[0].ID != courses[0].ID {
		t.Fatalf("peer courses = %+v, %v", peerCourses, err)
	}
	items, err := study.NewAssignmentService(tx, core.Access, core.Kernel.Notifier).List(ctx, student, courses[0].ID)
	if err != nil || len(items) != 2 {
		t.Fatalf("assignments = %+v, %v", items, err)
	}
	chats, err := chat.NewService(core.Kernel).List(ctx, student)
	if err != nil || len(chats) != 3 {
		t.Fatalf("student chats = %+v, %v", chats, err)
	}
	for _, c := range chats {
		if c.Request || c.Last == nil {
			t.Errorf("chat %q is a request or has no messages", c.Title)
		}
	}
	today := study.DayOf(time.Now())
	lessons, err := study.NewScheduleService(core.Kernel).Mine(ctx, student, today, today.AddDays(6))
	if err != nil || len(lessons) != 2 {
		t.Fatalf("week lessons = %+v, %v", lessons, err)
	}
}
