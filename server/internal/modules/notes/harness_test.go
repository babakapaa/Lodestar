package notes_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/files"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/live"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/core/tenant"
	"github.com/babakapaa/lodestar/server/internal/modules/notes"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/database/dbtest"
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

type world struct {
	t        *testing.T
	ctx      context.Context
	tenant   string
	kernel   module.Kernel
	notes    notes.Service
	users    identity.Store
	settings *fakeSettings
}

func newWorld(t *testing.T) *world {
	t.Helper()
	ctx := context.Background()
	pool := dbtest.New(t)
	sets := []database.Migrations{tenant.Migrations(), identity.Migrations(), live.Migrations(), files.Migrations(), notes.Module{}.Migrations()}
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
	settings := &fakeSettings{values: map[string]any{notes.QuotaSetting: 10240, notes.TrashDaysSetting: 30}}
	k := module.Kernel{Pool: pool, Tenants: tx, Settings: settings, Secrets: box, Outbound: &http.Client{Transport: fakeWeb{}}}
	return &world{t: t, ctx: ctx, tenant: created.ID, kernel: k, notes: notes.NewService(k), users: identity.NewStore(tx), settings: settings}
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

func (w *world) create(p auth.Principal, path string) notes.Note {
	w.t.Helper()
	n, _, err := w.notes.Create(w.ctx, p, notes.CreateInput{ID: noteID, Meta: notes.Meta{Path: path}})
	if err != nil {
		w.t.Fatalf("create %s: %v", path, err)
	}
	return n
}

func (w *world) push(p auth.Principal, key, data string) int64 {
	w.t.Helper()
	u, _, err := w.notes.Push(w.ctx, p, noteID, notes.UpdateInput{Key: key, Data: []byte(data)})
	if err != nil {
		w.t.Fatalf("push %s: %v", key, err)
	}
	return u.Seq
}

func (w *world) raw(query string, args ...any) []byte {
	w.t.Helper()
	var out []byte
	err := w.kernel.Tenants.Run(w.ctx, w.tenant, func(tx pgx.Tx) error {
		return tx.QueryRow(w.ctx, query, args...).Scan(&out)
	})
	if err != nil {
		w.t.Fatalf("raw: %v", err)
	}
	return out
}

func (w *world) serve(p auth.Principal, method, path, body string) *httptest.ResponseRecorder {
	w.t.Helper()
	mux := http.NewServeMux()
	notes.Routes{Service: w.notes}.Mount(mux)
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req = req.WithContext(auth.WithPrincipal(req.Context(), p))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}

type fakeWeb struct{}

func (fakeWeb) RoundTrip(r *http.Request) (*http.Response, error) {
	rec := httptest.NewRecorder()
	switch {
	case r.URL.Host == "api.crossref.org" && r.URL.Path == "/works/10.1000/xyz":
		rec.Header().Set("Content-Type", "application/json")
		_, _ = rec.WriteString(`{"message":{"type":"journal-article","title":["Строение клетки"],` +
			`"author":[{"given":"А. Б.","family":"Иванов"},{"given":"C.","family":"Smith"}],` +
			`"container-title":["Вестник биологии"],"volume":"12","issue":"3","page":"45-67",` +
			`"DOI":"10.1000/xyz","issued":{"date-parts":[[2021,5]]}}}`)
	case r.URL.Host == "journal.example.org":
		rec.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = rec.WriteString(`<html><head><title>Сайт</title>` +
			`<meta name="citation_title" content="Митоз и мейоз">` +
			`<meta name="citation_author" content="Ким, Е.">` +
			`<meta name="citation_author" content="Ли, А.">` +
			`<meta name="citation_journal_title" content="Биология">` +
			`<meta name="citation_publication_date" content="2019/03/01">` +
			`<meta name="citation_doi" content="https://doi.org/10.1/m">` +
			`</head><body>text</body></html>`)
	default:
		rec.WriteHeader(http.StatusNotFound)
	}
	return rec.Result(), nil
}
