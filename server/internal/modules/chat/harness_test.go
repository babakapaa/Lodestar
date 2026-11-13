package chat_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/audit"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/files"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/live"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/core/search"
	"github.com/babakapaa/lodestar/server/internal/core/structure"
	"github.com/babakapaa/lodestar/server/internal/core/tenant"
	"github.com/babakapaa/lodestar/server/internal/modules/chat"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/database/dbtest"
	"github.com/babakapaa/lodestar/server/internal/platform/secretbox"
)

type fakeSettings struct {
	retention int
}

func (f *fakeSettings) Value(context.Context, pgx.Tx, string) (json.RawMessage, error) {
	return json.Marshal(f.retention)
}

func (*fakeSettings) ModuleEnabled(context.Context, string, string) (bool, error) {
	return true, nil
}

type outbox struct {
	mu   sync.Mutex
	list []module.Notification
}

func (o *outbox) Notify(_ context.Context, _ pgx.Tx, _ string, n module.Notification) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.list = append(o.list, n)
	return nil
}

func (o *outbox) last() module.Notification {
	o.mu.Lock()
	defer o.mu.Unlock()
	if len(o.list) == 0 {
		return module.Notification{}
	}
	return o.list[len(o.list)-1]
}

type fakeRosters struct {
	mu   sync.Mutex
	list map[string]module.Roster
}

func (*fakeRosters) Kind() string {
	return "study.course"
}

func (f *fakeRosters) All(context.Context, pgx.Tx) ([]module.Roster, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make([]module.Roster, 0, len(f.list))
	for _, r := range f.list {
		out = append(out, r)
	}
	return out, nil
}

func (f *fakeRosters) One(_ context.Context, _ pgx.Tx, id string) (module.Roster, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	r, ok := f.list[id]
	if !ok {
		return module.Roster{}, module.ErrRosterNotFound
	}
	return r, nil
}

func (f *fakeRosters) set(r module.Roster) {
	f.mu.Lock()
	defer f.mu.Unlock()
	r.Kind = "study.course"
	f.list[r.ID] = r
}

type world struct {
	t        *testing.T
	ctx      context.Context
	tenant   string
	tx       database.TenantTx
	users    identity.Store
	access   access.Service
	chats    chat.Service
	messages chat.MessageService
	worker   chat.RosterWorker
	blocks   chat.BlockService
	reports  chat.ReportService
	purger   chat.Purger
	rosters  *fakeRosters
	sent     *outbox
	settings *fakeSettings
	pool     *pgxpool.Pool
	admin    auth.Principal
}

func newWorld(t *testing.T) *world {
	t.Helper()
	ctx := context.Background()
	pool := dbtest.New(t)
	sets := []database.Migrations{
		tenant.Migrations(), search.Migrations(), identity.Migrations(), audit.Migrations(), access.Migrations(), structure.Migrations(), live.Migrations(), files.Migrations(),
		chat.Module{}.Migrations(),
	}
	if err := database.NewMigrator(pool).Apply(ctx, sets...); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	created, err := tenant.NewStore(pool).Create(ctx, "kaznu", "КазНУ")
	if err != nil {
		t.Fatalf("tenant: %v", err)
	}
	catalog, err := access.NewCatalog(append(structure.Permissions(), chat.Module{}.Permissions()...)...)
	if err != nil {
		t.Fatalf("catalog: %v", err)
	}
	box, err := secretbox.New(base64.StdEncoding.EncodeToString(bytes.Repeat([]byte{7}, secretbox.KeySize)))
	if err != nil {
		t.Fatalf("box: %v", err)
	}
	tx := database.NewTenantTx(pool)
	acc := access.NewService(tx, catalog, structure.Tree{})
	rosters := &fakeRosters{list: map[string]module.Roster{}}
	sent := &outbox{}
	settings := &fakeSettings{}
	k := module.Kernel{
		Pool: pool, Tenants: tx, Access: acc, Notifier: sent, Settings: settings, Secrets: box,
		Rosters: module.Rosters{rosters},
	}
	w := &world{
		t: t, ctx: ctx, tenant: created.ID, tx: tx, users: identity.NewStore(tx), access: acc,
		chats: chat.NewService(k), messages: chat.NewMessageService(k), worker: chat.NewRosterWorker(k),
		blocks: chat.NewBlockService(k), reports: chat.NewReportService(k), purger: chat.NewPurger(k),
		rosters: rosters, sent: sent, settings: settings, pool: pool,
	}
	w.admin = w.user("admin")
	if _, err := access.NewBootstrap(tx).GrantTopRole(ctx, w.tenant, w.admin.UserID); err != nil {
		t.Fatalf("grant: %v", err)
	}
	return w
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

func (w *world) grant(p auth.Principal, perms ...access.Permission) {
	w.t.Helper()
	role, err := w.access.CreateRole(w.ctx, w.admin, access.RoleInput{Name: "r-" + p.UserID[:8], Permissions: perms})
	if err != nil {
		w.t.Fatalf("role: %v", err)
	}
	if _, err := w.access.Assign(w.ctx, w.admin, access.AssignInput{RoleID: role.ID, UserID: p.UserID, Scope: "tenant"}); err != nil {
		w.t.Fatalf("assign: %v", err)
	}
}

func (w *world) send(p auth.Principal, chatID, clientID, text string) chat.Message {
	w.t.Helper()
	m, _, err := w.messages.Send(w.ctx, p, chatID, chat.MessageInput{ClientID: clientID, Text: text})
	if err != nil {
		w.t.Fatalf("send %q: %v", text, err)
	}
	return m
}

func (w *world) sync() {
	w.t.Helper()
	if err := w.worker.SyncTenant(w.ctx, w.tenant); err != nil {
		w.t.Fatalf("sync: %v", err)
	}
}

func (w *world) events(p auth.Principal, kind string) int {
	w.t.Helper()
	var n int
	err := w.tx.Run(w.ctx, w.tenant, func(tx pgx.Tx) error {
		return tx.QueryRow(w.ctx, `SELECT count(*) FROM core.live_events WHERE user_id::text = $1 AND type = $2`, p.UserID, kind).Scan(&n)
	})
	if err != nil {
		w.t.Fatalf("events: %v", err)
	}
	return n
}

func (w *world) raw(query string) []byte {
	w.t.Helper()
	var out []byte
	err := w.tx.Run(w.ctx, w.tenant, func(tx pgx.Tx) error {
		return tx.QueryRow(w.ctx, query).Scan(&out)
	})
	if err != nil {
		w.t.Fatalf("raw: %v", err)
	}
	return out
}

func (w *world) call(p auth.Principal, method, path, body string) *httptest.ResponseRecorder {
	w.t.Helper()
	mux := http.NewServeMux()
	chat.Routes{Chats: w.chats, Messages: w.messages, Blocks: w.blocks, Reports: w.reports}.Mount(mux)
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req = req.WithContext(auth.WithPrincipal(req.Context(), p))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}

func (w *world) file(owner auth.Principal, name string) string {
	w.t.Helper()
	var id string
	err := w.tx.Run(w.ctx, w.tenant, func(tx pgx.Tx) error {
		f, err := files.FileStore{}.Create(w.ctx, tx, owner.UserID, files.FileInput{Name: name, ContentType: "text/plain", Size: 1})
		if err != nil {
			return err
		}
		id = f.ID
		_, err = files.FileStore{}.MarkReady(w.ctx, tx, f.ID)
		return err
	})
	if err != nil {
		w.t.Fatalf("file: %v", err)
	}
	return id
}

func (w *world) direct(a, b auth.Principal) chat.ChatDetails {
	w.t.Helper()
	d, _, err := w.chats.Direct(w.ctx, a, chat.DirectInput{UserID: b.UserID})
	if err != nil {
		w.t.Fatalf("direct: %v", err)
	}
	return d
}
