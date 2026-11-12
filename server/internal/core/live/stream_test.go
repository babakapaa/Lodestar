package live_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/live"
	"github.com/babakapaa/lodestar/server/internal/core/tenant"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/database/dbtest"
)

type harness struct {
	t      *testing.T
	ctx    context.Context
	tx     database.TenantTx
	tenant string
	user   string
	token  string
	server *httptest.Server
}

func newHarness(t *testing.T) *harness {
	t.Helper()
	pool := dbtest.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	sets := []database.Migrations{tenant.Migrations(), identity.Migrations(), auth.Migrations(), live.Migrations()}
	if err := database.NewMigrator(pool).Apply(ctx, sets...); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	created, err := tenant.NewStore(pool).Create(ctx, "kaznu", "КазНУ")
	if err != nil {
		t.Fatalf("tenant: %v", err)
	}
	tx := database.NewTenantTx(pool)
	u, err := identity.NewStore(tx).UpsertExternal(ctx, created.ID, identity.ExternalIdentity{Issuer: "https://idp", Subject: "s"})
	if err != nil {
		t.Fatalf("user: %v", err)
	}
	sessions := auth.NewSessionStore(tx, auth.Lifetimes{Access: time.Minute, Refresh: time.Hour})
	pair, err := sessions.Start(ctx, created.ID, u.ID, "test")
	if err != nil {
		t.Fatalf("session: %v", err)
	}
	hub := live.NewHub()
	go live.NewListener(pool, hub).Run(ctx)
	mux := http.NewServeMux()
	live.Routes{Hub: hub, Tenants: tx, Authenticator: auth.NewAuthenticator(sessions)}.Mount(mux)
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	return &harness{t: t, ctx: ctx, tx: tx, tenant: created.ID, user: u.ID, token: pair.AccessToken, server: server}
}

func (h *harness) publish(kind string) {
	h.t.Helper()
	err := h.tx.Run(h.ctx, h.tenant, func(tx pgx.Tx) error {
		return live.Store{}.Publish(h.ctx, tx, h.tenant, []string{h.user}, kind, map[string]string{"k": kind})
	})
	if err != nil {
		h.t.Fatalf("publish: %v", err)
	}
}

func (h *harness) dial(after int64) *websocket.Conn {
	h.t.Helper()
	url := "ws" + strings.TrimPrefix(h.server.URL, "http") + "/v1/live?after=" + strconv.FormatInt(after, 10) + "&access_token=" + h.token
	conn, resp, err := websocket.Dial(h.ctx, url, nil)
	if err != nil {
		h.t.Fatalf("dial: %v", err)
	}
	if resp.Body != nil {
		_ = resp.Body.Close()
	}
	h.t.Cleanup(func() { _ = conn.CloseNow() })
	return conn
}

func read(t *testing.T, conn *websocket.Conn) live.Event {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var e live.Event
	if err := wsjson.Read(ctx, conn, &e); err != nil {
		t.Fatalf("read: %v", err)
	}
	return e
}

func TestStreamDeliversLiveEvents(t *testing.T) {
	h := newHarness(t)
	conn := h.dial(0)
	time.Sleep(200 * time.Millisecond)

	h.publish("test.first")
	h.publish("test.second")

	if e := read(t, conn); e.Type != "test.first" {
		t.Fatalf("first = %+v", e)
	}
	if e := read(t, conn); e.Type != "test.second" {
		t.Fatalf("second = %+v", e)
	}
}

func TestStreamReplaysMissedEvents(t *testing.T) {
	h := newHarness(t)
	h.publish("test.one")
	h.publish("test.two")

	first := read(t, h.dial(0))
	conn := h.dial(first.ID)
	if e := read(t, conn); e.Type != "test.two" {
		t.Fatalf("replay after %d = %+v", first.ID, e)
	}
}

func TestStreamRequiresToken(t *testing.T) {
	h := newHarness(t)
	url := "ws" + strings.TrimPrefix(h.server.URL, "http") + "/v1/live"

	_, resp, err := websocket.Dial(h.ctx, url, nil)
	if resp != nil {
		defer func() { _ = resp.Body.Close() }()
	}
	if err == nil {
		t.Fatal("dial without token succeeded")
	}
	if resp == nil || resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("response = %+v", resp)
	}
}

func TestEventsEndpointPages(t *testing.T) {
	h := newHarness(t)
	h.publish("test.one")

	req, _ := http.NewRequest(http.MethodGet, h.server.URL+"/v1/events?after=0", nil)
	req.Header.Set("Authorization", "Bearer "+h.token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("events: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", resp.StatusCode)
	}
}
