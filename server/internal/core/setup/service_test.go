package setup_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/audit"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/core/schema"
	"github.com/babakapaa/lodestar/server/internal/core/search"
	"github.com/babakapaa/lodestar/server/internal/core/settings"
	"github.com/babakapaa/lodestar/server/internal/core/setup"
	"github.com/babakapaa/lodestar/server/internal/core/structure"
	"github.com/babakapaa/lodestar/server/internal/core/tenant"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/database/dbtest"
)

type lonelyCourses struct {
	count *int
}

func (l lonelyCourses) Problems(context.Context, pgx.Tx, string) ([]module.Problem, error) {
	if *l.count == 0 {
		return nil, nil
	}
	return []module.Problem{{Section: setup.SectionTeaching, Code: "courses_without_teacher", Count: *l.count}}, nil
}

type world struct {
	t         *testing.T
	ctx       context.Context
	tenant    string
	users     identity.Store
	access    access.Service
	structure structure.Service
	settings  settings.Service
	setup     setup.Service
	admin     auth.Principal
	lonely    int
	count     int
	pool      *pgxpool.Pool
	sent      *sentNotifications
	pilot     *fakePilot
}

func newWorld(t *testing.T) *world {
	t.Helper()
	ctx := context.Background()
	pool := dbtest.New(t)
	sets := []database.Migrations{
		tenant.Migrations(), search.Migrations(), identity.Migrations(), audit.Migrations(), access.Migrations(),
		structure.Migrations(), settings.Migrations(), setup.Migrations(),
	}
	if err := database.NewMigrator(pool).Apply(ctx, sets...); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	created, err := tenant.NewStore(pool).Create(ctx, "kaznu", "КазНУ")
	if err != nil {
		t.Fatalf("tenant: %v", err)
	}
	perms := append(append(structure.Permissions(), settings.Permissions()...), setup.Permissions()...)
	catalog, err := access.NewCatalog(perms...)
	if err != nil {
		t.Fatalf("catalog: %v", err)
	}
	defs, err := schema.NewCatalog(
		schema.Definition{Key: "core.self_deletion", Type: schema.Bool, Default: true},
		schema.Definition{Key: "core.time_zone", Type: schema.TimeZone, Default: "UTC"},
	)
	if err != nil {
		t.Fatalf("schema: %v", err)
	}
	registry, err := module.NewRegistry()
	if err != nil {
		t.Fatalf("registry: %v", err)
	}
	tx := database.NewTenantTx(pool)
	acc := access.NewService(tx, catalog, structure.Tree{})
	w := &world{
		t: t, ctx: ctx, tenant: created.ID, users: identity.NewStore(tx), access: acc, structure: structure.NewService(tx, acc),
		pool: pool, sent: &sentNotifications{}, pilot: &fakePilot{members: map[string][]string{}},
	}
	w.settings = settings.NewService(tx, acc, defs, registry)
	w.setup = setup.NewService(tx, acc, w.settings, w.sent, []module.PilotSource{w.pilot}, lonelyCourses{count: &w.lonely})
	w.admin = w.user()
	if _, err := access.NewBootstrap(tx).GrantTopRole(ctx, w.tenant, w.admin.UserID); err != nil {
		t.Fatalf("grant: %v", err)
	}
	return w
}

func (w *world) user() auth.Principal {
	w.t.Helper()
	w.count++
	subject := fmt.Sprintf("u%d", w.count)
	u, err := w.users.UpsertExternal(w.ctx, w.tenant, identity.ExternalIdentity{Issuer: "https://idp", Subject: subject, Email: subject + "@kaznu.kz"})
	if err != nil {
		w.t.Fatalf("user: %v", err)
	}
	return auth.Principal{TenantID: w.tenant, UserID: u.ID}
}

func (w *world) grant(p auth.Principal, perms ...access.Permission) {
	w.t.Helper()
	role, err := w.access.CreateRole(w.ctx, w.admin, access.RoleInput{Name: fmt.Sprintf("r%d", w.count), Permissions: perms})
	if err != nil {
		w.t.Fatalf("role: %v", err)
	}
	if _, err := w.access.Assign(w.ctx, w.admin, access.AssignInput{RoleID: role.ID, UserID: p.UserID, Scope: "tenant"}); err != nil {
		w.t.Fatalf("assign: %v", err)
	}
}

func codes(problems []module.Problem) []string {
	out := []string{}
	for _, p := range problems {
		out = append(out, p.Section+"/"+p.Code)
	}
	return out
}

func TestLaunchedUniversityIsOpen(t *testing.T) {
	w := newWorld(t)
	student := w.user()
	if ok, err := w.setup.Allowed(w.ctx, student); err != nil || !ok {
		t.Fatalf("allowed = %v, %v", ok, err)
	}
	if _, err := w.setup.Overview(w.ctx, w.admin); !errors.Is(err, setup.ErrNoSetup) {
		t.Fatalf("overview: %v", err)
	}
	view, err := w.setup.Status(w.ctx, student)
	if err != nil || !view.Launched {
		t.Fatalf("status = %+v, %v", view, err)
	}
}

func TestWizardFromDelegationToLaunch(t *testing.T) {
	w := newWorld(t)
	if err := w.setup.Begin(w.ctx, w.tenant); err != nil {
		t.Fatalf("begin: %v", err)
	}
	student, it, lawyer, rector := w.user(), w.user(), w.user(), w.user()
	w.grant(lawyer, setup.Legal)
	w.grant(rector, setup.Sign)
	w.lonely = 2

	if ok, err := w.setup.Allowed(w.ctx, student); err != nil || ok {
		t.Fatalf("student before launch = %v, %v", ok, err)
	}
	if _, err := w.setup.Overview(w.ctx, student); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("student overview: %v", err)
	}
	if _, err := w.setup.Assign(w.ctx, w.admin, setup.SectionLegal, setup.AssignInput{UserID: &it.UserID}); !errors.Is(err, setup.ErrInvalidAssignee) {
		t.Fatalf("legal to it: %v", err)
	}
	due := "2026-12-01"
	sec, err := w.setup.Assign(w.ctx, w.admin, setup.SectionUniversity, setup.AssignInput{UserID: &it.UserID, DueOn: &due})
	if err != nil || sec.AssigneeID == nil || *sec.DueOn != due || len(sec.Settings) != 1 || sec.Settings[0] != "core.time_zone" {
		t.Fatalf("assign = %+v, %v", sec, err)
	}
	if ok, err := w.setup.Allowed(w.ctx, it); err != nil || !ok {
		t.Fatalf("assignee allowed = %v, %v", ok, err)
	}
	it2, err := w.setup.Overview(w.ctx, it)
	if err != nil || len(it2.Sections) != 1 || it2.Sections[0].Key != setup.SectionUniversity || it2.Problems != nil || len(it2.Legal) != 0 {
		t.Fatalf("assignee overview = %+v, %v", it2, err)
	}
	if err := w.setup.SetSetting(w.ctx, it, setup.SectionUniversity, "core.self_deletion", json.RawMessage(`false`)); !errors.Is(err, setup.ErrWrongSection) {
		t.Fatalf("other section setting: %v", err)
	}
	if err := w.setup.SetSetting(w.ctx, it, setup.SectionPrivacy, "core.self_deletion", json.RawMessage(`false`)); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("not assigned section: %v", err)
	}
	if err := w.setup.SetSetting(w.ctx, it, setup.SectionUniversity, "core.time_zone", json.RawMessage(`"Asia/Almaty"`)); err != nil {
		t.Fatalf("set: %v", err)
	}
	if sub, err := w.setup.Submit(w.ctx, it, setup.SectionUniversity); err != nil || sub.Status != setup.InReview {
		t.Fatalf("submit = %+v, %v", sub, err)
	}
	if _, err := w.setup.Submit(w.ctx, it, setup.SectionUniversity); !errors.Is(err, setup.ErrInvalidStatus) {
		t.Fatalf("submit twice: %v", err)
	}
	if _, err := w.setup.Review(w.ctx, it, setup.SectionUniversity, setup.ReviewInput{Approve: true}); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("assignee approves: %v", err)
	}
	back, err := w.setup.Review(w.ctx, w.admin, setup.SectionUniversity, setup.ReviewInput{Note: "Проверьте часовой пояс"})
	if err != nil || back.Status != setup.InProgress || back.Note == "" {
		t.Fatalf("return = %+v, %v", back, err)
	}

	if _, err := w.setup.Answer(w.ctx, student, "storage", setup.AnswerInput{Answer: "a", Basis: "b"}); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("student answers: %v", err)
	}
	if _, err := w.setup.Answer(w.ctx, lawyer, "storage", setup.AnswerInput{Answer: "Серверы вуза"}); !errors.Is(err, setup.ErrInvalidAnswer) {
		t.Fatalf("no basis: %v", err)
	}
	for _, q := range setup.LegalQuestions {
		if _, err := w.setup.Answer(w.ctx, lawyer, q.Key, setup.AnswerInput{Answer: "Решено", Basis: "Приказ ректора №1"}); err != nil {
			t.Fatalf("answer %s: %v", q.Key, err)
		}
	}
	if _, err := w.setup.Sign(w.ctx, rector); !errors.Is(err, setup.ErrLegalOpen) {
		t.Fatalf("sign before approval: %v", err)
	}
	if _, err := w.setup.ApproveAnswer(w.ctx, w.user(), "storage"); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("stranger approves: %v", err)
	}
	for _, q := range setup.LegalQuestions {
		if _, err := w.setup.ApproveAnswer(w.ctx, lawyer, q.Key); err != nil {
			t.Fatalf("approve %s: %v", q.Key, err)
		}
	}
	if _, err := w.setup.Sign(w.ctx, w.user()); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("stranger signs: %v", err)
	}
	signed, err := w.setup.Sign(w.ctx, rector)
	if err != nil || signed.ProtocolVersion != 1 || signed.SignedAt == nil {
		t.Fatalf("sign = %+v, %v", signed, err)
	}
	protocol, err := w.setup.Protocol(w.ctx, lawyer)
	if text := protocol.Markdown(); err != nil || !strings.Contains(text, "Основание: Приказ ректора №1") || !strings.Contains(text, "версия 1") {
		t.Fatalf("protocol = %s, %v", text, err)
	}
	if doc, err := protocol.PDF(); err != nil || !bytes.HasPrefix(doc, []byte("%PDF-")) {
		t.Fatalf("pdf = %d bytes, %v", len(doc), err)
	}
	if _, err := w.setup.Answer(w.ctx, lawyer, "recording", setup.AnswerInput{Answer: "Только с согласия", Basis: "Приказ №2"}); err != nil {
		t.Fatalf("change answer: %v", err)
	}
	o, err := w.setup.Overview(w.ctx, w.admin)
	if err != nil || o.SignedAt != nil || o.ProtocolVersion != 1 {
		t.Fatalf("after change = %+v, %v", o, err)
	}
	want := "deploy/section_not_approved university/section_not_approved people/section_not_approved admins/section_not_approved " +
		"legal/section_not_approved legal/legal_open launch/protocol_unsigned university/no_units teaching/courses_without_teacher"
	if got := strings.Join(codes(o.Problems), " "); got != want {
		t.Fatalf("problems = %s", got)
	}
	if _, err := w.setup.Launch(w.ctx, w.admin, setup.LaunchInput{}); !errors.Is(err, setup.ErrNotReady) {
		t.Fatalf("launch early: %v", err)
	}

	if _, err := w.setup.ApproveAnswer(w.ctx, lawyer, "recording"); err != nil {
		t.Fatalf("approve again: %v", err)
	}
	if _, err := w.setup.Sign(w.ctx, rector); err != nil {
		t.Fatalf("sign again: %v", err)
	}
	for _, key := range []string{setup.SectionDeploy, setup.SectionUniversity, setup.SectionPeople, setup.SectionAdmins, setup.SectionLegal} {
		if _, err := w.setup.Review(w.ctx, w.admin, key, setup.ReviewInput{Approve: true}); err != nil {
			t.Fatalf("approve %s: %v", key, err)
		}
	}
	if _, err := w.structure.Create(w.ctx, w.admin, structure.UnitInput{Kind: "faculty", Name: "Биофак"}); err != nil {
		t.Fatalf("unit: %v", err)
	}
	w.lonely = 0
	launched, err := w.setup.Launch(w.ctx, w.admin, setup.LaunchInput{})
	if err != nil || !launched.Launched || launched.ProtocolVersion != 2 || len(launched.Problems) != 0 {
		t.Fatalf("launch = %+v, %v", launched, err)
	}
	for _, s := range launched.Sections {
		if s.AssigneeID != nil {
			t.Fatalf("assignee kept = %+v", s)
		}
	}
	if ok, err := w.setup.Allowed(w.ctx, student); err != nil || !ok {
		t.Fatalf("student after launch = %v, %v", ok, err)
	}
}
