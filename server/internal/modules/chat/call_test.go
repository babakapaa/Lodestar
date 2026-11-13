package chat_test

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"

	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/modules/chat"
)

const callSecret = "0123456789abcdef0123456789abcdef"

func hook(t *testing.T, w *world, body string) int {
	t.Helper()
	sum := sha256.Sum256([]byte(body))
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS256, Key: []byte(callSecret)}, nil)
	if err != nil {
		t.Fatalf("signer: %v", err)
	}
	auth, err := jwt.Signed(signer).Claims(jwt.Claims{Issuer: "key", Expiry: jwt.NewNumericDate(time.Now().Add(time.Minute))}).
		Claims(map[string]any{"sha256": base64.StdEncoding.EncodeToString(sum[:])}).Serialize()
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	mux := http.NewServeMux()
	chat.Routes{Calls: w.calls, Recordings: w.records}.MountOpen(mux)
	req := httptest.NewRequest(http.MethodPost, "/v1/open/chat/livekit", strings.NewReader(body))
	req.Header.Set("Authorization", auth)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec.Code
}

func TestGroupCallFromRingToEnd(t *testing.T) {
	w := newWorld(t)
	owner, member, stranger := w.user("owner"), w.user("member"), w.user("stranger")
	g, _, err := w.chats.Create(w.ctx, owner, chat.ChatInput{Kind: chat.KindGroup, Title: "Команда", MemberIDs: []string{member.UserID, owner.UserID}})
	if err != nil {
		t.Fatalf("group: %v", err)
	}
	if _, err := w.calls.Active(w.ctx, owner, g.ID); !errors.Is(err, chat.ErrCallNotFound) {
		t.Fatalf("no call yet: %v", err)
	}
	if _, err := w.calls.Join(w.ctx, stranger, g.ID); !errors.Is(err, chat.ErrChatNotFound) {
		t.Fatalf("stranger joins: %v", err)
	}
	first, err := w.calls.Join(w.ctx, owner, g.ID)
	if err != nil || first.Token == "" || first.URL != "wss://calls.kaznu.kz" || !first.ScreenShare || first.Room != chat.RoomName(w.tenant, first.Call.ID) {
		t.Fatalf("owner joins = %+v, %v", first, err)
	}
	if n := w.sent.last(); n.Type != chat.NotifyCall || len(n.UserIDs) != 1 || n.UserIDs[0] != member.UserID {
		t.Fatalf("ring = %+v", n)
	}
	second, err := w.calls.Join(w.ctx, member, g.ID)
	if err != nil || second.Call.ID != first.Call.ID || second.ScreenShare || len(second.Call.People) != 2 {
		t.Fatalf("member joins = %+v, %v", second, err)
	}
	if err := w.calls.Leave(w.ctx, stranger, first.Call.ID); !errors.Is(err, chat.ErrCallNotFound) {
		t.Fatalf("stranger leaves: %v", err)
	}
	if err := w.calls.Leave(w.ctx, owner, first.Call.ID); err != nil {
		t.Fatalf("owner leaves: %v", err)
	}
	active, err := w.calls.Active(w.ctx, member, g.ID)
	if err != nil || len(active.People) != 1 || active.People[0] != member.UserID {
		t.Fatalf("after leave = %+v, %v", active, err)
	}
	body := `{"event":"participant_left","room":{"name":"` + first.Room + `"},"participant":{"identity":"` + member.UserID + `"}}`
	if code := hook(t, w, body); code != http.StatusNoContent {
		t.Fatalf("hook = %d", code)
	}
	if _, err := w.calls.Active(w.ctx, member, g.ID); !errors.Is(err, chat.ErrCallNotFound) {
		t.Fatalf("call ended: %v", err)
	}
	again, err := w.calls.Join(w.ctx, member, g.ID)
	if err != nil || again.Call.ID == first.Call.ID {
		t.Fatalf("new call = %+v, %v", again, err)
	}
	if code := hook(t, w, `{"event":"room_finished","room":{"name":"`+again.Room+`"}}`); code != http.StatusNoContent {
		t.Fatalf("finish hook = %d", code)
	}
	if _, err := w.calls.Active(w.ctx, member, g.ID); !errors.Is(err, chat.ErrCallNotFound) {
		t.Fatalf("room finished: %v", err)
	}

	w.settings.values[chat.CallsSetting] = false
	if _, err := w.calls.Join(w.ctx, owner, g.ID); !errors.Is(err, chat.ErrCallsOff) {
		t.Fatalf("calls off: %v", err)
	}
}

func TestCallHookNeedsSignature(t *testing.T) {
	w := newWorld(t)
	mux := http.NewServeMux()
	chat.Routes{Calls: w.calls, Recordings: w.records}.MountOpen(mux)
	req := httptest.NewRequest(http.MethodPost, "/v1/open/chat/livekit", strings.NewReader(`{}`))
	req.Header.Set("Authorization", "nope")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("unsigned = %d", rec.Code)
	}
}

func TestCallsWithoutServer(t *testing.T) {
	w := newWorld(t)
	owner := w.user("owner")
	g, _, err := w.chats.Create(w.ctx, owner, chat.ChatInput{Kind: chat.KindGroup, Title: "Команда", MemberIDs: []string{owner.UserID}})
	if err != nil {
		t.Fatalf("group: %v", err)
	}
	calls := chat.NewCallService(module.Kernel{Tenants: w.tx})
	if _, err := calls.Join(w.ctx, owner, g.ID); !errors.Is(err, chat.ErrCallsUnavailable) {
		t.Fatalf("no server: %v", err)
	}
}
