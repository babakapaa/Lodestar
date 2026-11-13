package navigator_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/core/tenant"
	"github.com/babakapaa/lodestar/server/internal/modules/navigator"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/database/dbtest"
	"github.com/babakapaa/lodestar/server/internal/platform/llm"
	"github.com/babakapaa/lodestar/server/internal/platform/secretbox"
)

type fakeSettings struct {
	values map[string]any
}

func (f *fakeSettings) Value(_ context.Context, _ pgx.Tx, key string) (json.RawMessage, error) {
	return json.Marshal(f.values[key])
}

func (f *fakeSettings) ModuleEnabled(context.Context, string, string) (bool, error) {
	return true, nil
}

const assignmentID = "5a6b7c8d-9e0f-4a1b-8c2d-3e4f5a6b7c8d"

type fakeBriefer struct {
	mode     string
	students []string
}

func (fakeBriefer) Kind() string {
	return "assignment"
}

func (b *fakeBriefer) Brief(_ context.Context, _ pgx.Tx, p auth.Principal, id string) (module.Brief, error) {
	if id != assignmentID || !slices.Contains(b.students, p.UserID) {
		return module.Brief{}, module.ErrBriefNotFound
	}
	due := time.Date(2030, 12, 1, 18, 0, 0, 0, time.UTC)
	return module.Brief{
		Kind: "assignment", ID: id, Title: "Реферат о клетке", Description: "Опишите строение клетки",
		Course: "Биология", DueAt: &due, Mode: b.mode,
	}, nil
}

type fakeModel struct {
	mu        sync.Mutex
	intent    string
	distress  bool
	reply     string
	fail      bool
	requests  []llm.Request
	completes int
}

func (m *fakeModel) Complete(_ context.Context, req llm.Request) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.requests = append(m.requests, req)
	if req.JSON {
		raw, _ := json.Marshal(map[string]any{"intent": m.intent, "distress": m.distress})
		return "```json\n" + string(raw) + "\n```", nil
	}
	m.completes++
	if m.fail {
		return "", llm.ErrModelFailed
	}
	return m.reply, nil
}

func (m *fakeModel) Stream(ctx context.Context, req llm.Request, delta func(string) error) (string, error) {
	full, err := m.Complete(ctx, req)
	if err != nil {
		return "", err
	}
	for _, piece := range strings.SplitAfter(full, " ") {
		if err := delta(piece); err != nil {
			return "", err
		}
	}
	return full, nil
}

func (m *fakeModel) last() llm.Request {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.requests[len(m.requests)-1]
}

type world struct {
	t         *testing.T
	ctx       context.Context
	tenant    string
	kernel    module.Kernel
	navigator navigator.Service
	users     identity.Store
	model     *fakeModel
	briefer   *fakeBriefer
}

func newWorld(t *testing.T) *world {
	t.Helper()
	ctx := context.Background()
	pool := dbtest.New(t)
	sets := []database.Migrations{tenant.Migrations(), identity.Migrations(), navigator.Module{}.Migrations()}
	if err := database.NewMigrator(pool).Apply(ctx, sets...); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	created, err := tenant.NewStore(pool).Create(ctx, "kaznu", "КазНУ")
	if err != nil {
		t.Fatalf("tenant: %v", err)
	}
	box, err := secretbox.New(base64.StdEncoding.EncodeToString(bytes.Repeat([]byte{7}, secretbox.KeySize)))
	if err != nil {
		t.Fatalf("box: %v", err)
	}
	tx := database.NewTenantTx(pool)
	settings := &fakeSettings{values: map[string]any{
		navigator.SupportSetting: "Психологическая служба: +7 727 000 00 00",
		navigator.MemorySetting:  30,
	}}
	model := &fakeModel{intent: "understand", reply: "Что ты уже знаешь о клетке?"}
	briefer := &fakeBriefer{mode: module.ModeStudy}
	k := module.Kernel{Pool: pool, Tenants: tx, Settings: settings, Secrets: box, Model: model, Briefs: module.Briefs{briefer}}
	return &world{
		t: t, ctx: ctx, tenant: created.ID, kernel: k, navigator: navigator.NewService(k),
		users: identity.NewStore(tx), model: model, briefer: briefer,
	}
}

func (w *world) user(subject string) auth.Principal {
	w.t.Helper()
	u, err := w.users.UpsertExternal(w.ctx, w.tenant, identity.ExternalIdentity{
		Issuer: "https://idp", Subject: subject, Email: subject + "@kaznu.kz", DisplayName: subject,
	})
	if err != nil {
		w.t.Fatalf("user: %v", err)
	}
	return auth.Principal{TenantID: w.tenant, UserID: u.ID}
}

func (w *world) enrol(mode string, student auth.Principal) {
	w.briefer.mode = mode
	w.briefer.students = append(w.briefer.students, student.UserID)
}

func (w *world) post(p auth.Principal, path, body string, header map[string]string) *httptest.ResponseRecorder {
	w.t.Helper()
	mux := http.NewServeMux()
	navigator.Routes{Service: w.navigator}.Mount(mux)
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	for k, v := range header {
		req.Header.Set(k, v)
	}
	req = req.WithContext(auth.WithPrincipal(req.Context(), p))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}

func (w *world) raw(query string) []byte {
	w.t.Helper()
	var out []byte
	err := w.kernel.Tenants.Run(w.ctx, w.tenant, func(tx pgx.Tx) error {
		return tx.QueryRow(w.ctx, query).Scan(&out)
	})
	if err != nil {
		w.t.Fatalf("raw: %v", err)
	}
	return out
}
