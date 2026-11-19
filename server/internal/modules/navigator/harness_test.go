package navigator_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"hash/fnv"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"
	"unicode"

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

type fakePreferences struct {
	off []string
}

func (f *fakePreferences) Value(_ context.Context, _ pgx.Tx, userID, _ string) (json.RawMessage, error) {
	return json.Marshal(!slices.Contains(f.off, userID))
}

const (
	assignmentID = "5a6b7c8d-9e0f-4a1b-8c2d-3e4f5a6b7c8d"
	courseID     = "1c2d3e4f-5a6b-4c7d-8e9f-0a1b2c3d4e5f"
)

type course struct {
	mode     string
	students []string
	teachers []string
}

type fakeBriefer struct {
	kind   string
	course *course
}

func (b fakeBriefer) Kind() string {
	return b.kind
}

func (b fakeBriefer) Brief(_ context.Context, _ pgx.Tx, p auth.Principal, id string) (module.Brief, error) {
	teaches := slices.Contains(b.course.teachers, p.UserID)
	if !teaches && !slices.Contains(b.course.students, p.UserID) {
		return module.Brief{}, module.ErrBriefNotFound
	}
	if b.kind == "course" && id == courseID {
		return module.Brief{
			Kind: "course", ID: id, Title: "Биология", Course: "Биология", CourseID: courseID, Teaches: teaches, Mode: b.course.mode,
		}, nil
	}
	if b.kind != "assignment" || id != assignmentID {
		return module.Brief{}, module.ErrBriefNotFound
	}
	due := time.Date(2030, 12, 1, 18, 0, 0, 0, time.UTC)
	return module.Brief{
		Kind: "assignment", ID: id, Title: "Реферат о клетке", Description: "Опишите строение клетки",
		Course: "Биология", CourseID: courseID, Teaches: teaches, DueAt: &due, Mode: b.course.mode,
	}, nil
}

type fakeEmbedder struct{}

func (fakeEmbedder) Name() string {
	return "fake-embed"
}

func (fakeEmbedder) Embed(_ context.Context, texts []string) ([][]float32, error) {
	out := make([][]float32, len(texts))
	for i, text := range texts {
		v := make([]float32, 32)
		v[0] = 0.01
		for _, word := range strings.FieldsFunc(strings.ToLower(text), func(r rune) bool { return !unicode.IsLetter(r) }) {
			if len([]rune(word)) < 4 {
				continue
			}
			h := fnv.New32a()
			_, _ = h.Write([]byte(string([]rune(word)[:4])))
			v[h.Sum32()%32]++
		}
		out[i] = v
	}
	return out, nil
}

type fakeModel struct {
	mu          sync.Mutex
	intent      string
	distress    bool
	reply       string
	fail        bool
	unsupported []int
	leak        bool
	requests    []llm.Request
	completes   int
}

func (m *fakeModel) Complete(_ context.Context, req llm.Request) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.requests = append(m.requests, req)
	if req.JSON && strings.Contains(req.Messages[0].Content, "проверяешь") {
		raw, _ := json.Marshal(map[string]any{"unsupported": m.unsupported, "leak": m.leak})
		return string(raw), nil
	}
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

func (m *fakeModel) verifier() llm.Request {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i := len(m.requests) - 1; i >= 0; i-- {
		if m.requests[i].JSON && strings.Contains(m.requests[i].Messages[0].Content, "проверяешь") {
			return m.requests[i]
		}
	}
	return llm.Request{}
}

func (m *fakeModel) last() llm.Request {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i := len(m.requests) - 1; i >= 0; i-- {
		if !m.requests[i].JSON {
			return m.requests[i]
		}
	}
	return llm.Request{}
}

type world struct {
	t         *testing.T
	ctx       context.Context
	tenant    string
	kernel    module.Kernel
	navigator navigator.Service
	users     identity.Store
	materials navigator.MaterialService
	topics    navigator.TopicService
	model     *fakeModel
	course    *course
	prefs     *fakePreferences
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
	c := &course{mode: module.ModeStudy}
	prefs := &fakePreferences{}
	k := module.Kernel{
		Pool: pool, Tenants: tx, Settings: settings, Preferences: prefs, Secrets: box, Model: model, Embedder: fakeEmbedder{},
		Briefs: module.Briefs{fakeBriefer{kind: "assignment", course: c}, fakeBriefer{kind: "course", course: c}},
	}
	return &world{
		t: t, ctx: ctx, tenant: created.ID, kernel: k, navigator: navigator.NewService(k),
		materials: navigator.NewMaterialService(k), topics: navigator.NewTopicService(k), users: identity.NewStore(tx), model: model, course: c, prefs: prefs,
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
	w.course.mode = mode
	w.course.students = append(w.course.students, student.UserID)
}

func (w *world) teach(teacher auth.Principal) {
	w.course.teachers = append(w.course.teachers, teacher.UserID)
}

func (w *world) get(p auth.Principal, path string) *httptest.ResponseRecorder {
	w.t.Helper()
	return w.serve(p, httptest.NewRequest(http.MethodGet, path, nil))
}

func (w *world) post(p auth.Principal, path, body string, header map[string]string) *httptest.ResponseRecorder {
	w.t.Helper()
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	for k, v := range header {
		req.Header.Set(k, v)
	}
	return w.serve(p, req)
}

func (w *world) serve(p auth.Principal, req *http.Request) *httptest.ResponseRecorder {
	w.t.Helper()
	mux := http.NewServeMux()
	navigator.Routes{Service: w.navigator, Materials: w.materials, Topics: w.topics}.Mount(mux)
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
