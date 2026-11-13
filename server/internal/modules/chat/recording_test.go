package chat_test

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/modules/chat"
)

func TestRecordingFollowsUniversityRules(t *testing.T) {
	w := newWorld(t)
	owner, member := w.user("owner"), w.user("member")
	g, _, err := w.chats.Create(w.ctx, owner, chat.ChatInput{Kind: chat.KindGroup, Title: "Команда", MemberIDs: []string{member.UserID, owner.UserID}})
	if err != nil {
		t.Fatalf("group: %v", err)
	}
	first, err := w.calls.Join(w.ctx, owner, g.ID, chat.JoinInput{})
	if err != nil {
		t.Fatalf("join: %v", err)
	}
	if _, err := w.calls.Join(w.ctx, member, g.ID, chat.JoinInput{}); err != nil {
		t.Fatalf("member joins: %v", err)
	}
	callID := first.Call.ID
	if first.Record {
		t.Fatal("record allowed while off")
	}
	if _, err := w.records.Start(w.ctx, owner, callID, chat.RecordingInput{}); !errors.Is(err, chat.ErrRecordingOff) {
		t.Fatalf("off by default: %v", err)
	}
	w.settings.values[chat.RecordingSetting] = chat.RecordAdmins
	if again, err := w.calls.Join(w.ctx, owner, g.ID, chat.JoinInput{}); err != nil || !again.Record {
		t.Fatalf("owner may record = %+v, %v", again, err)
	}
	if again, err := w.calls.Join(w.ctx, member, g.ID, chat.JoinInput{}); err != nil || again.Record {
		t.Fatalf("member may not record = %+v, %v", again, err)
	}
	if _, err := w.records.Start(w.ctx, member, callID, chat.RecordingInput{}); !errors.Is(err, chat.ErrRecordingOff) {
		t.Fatalf("member records: %v", err)
	}
	call, err := w.records.Start(w.ctx, owner, callID, chat.RecordingInput{})
	if err != nil || !call.Recording {
		t.Fatalf("owner records = %+v, %v", call, err)
	}
	if _, err := w.records.Start(w.ctx, owner, callID, chat.RecordingInput{}); !errors.Is(err, chat.ErrRecordingState) {
		t.Fatalf("twice: %v", err)
	}
	seen, err := w.calls.Active(w.ctx, member, g.ID)
	if err != nil || !seen.Recording {
		t.Fatalf("member sees recording = %+v, %v", seen, err)
	}
	key := w.egress.key(0)
	if !strings.HasPrefix(key, w.tenant+"/") {
		t.Fatalf("key = %s", key)
	}
	if _, err := w.records.Stop(w.ctx, member, callID); !errors.Is(err, chat.ErrRecordingOff) {
		t.Fatalf("member stops: %v", err)
	}
	if call, err = w.records.Stop(w.ctx, owner, callID); err != nil || call.Recording || len(w.egress.stopped) != 1 {
		t.Fatalf("stop = %+v, %v", call, err)
	}

	room := chat.RoomName(w.tenant, callID)
	body := `{"event":"egress_ended","egressInfo":{"egressId":"EG_1","roomName":"` + room + `","status":"EGRESS_COMPLETE","fileResults":[{"filename":"` + key + `","size":"4096"}]}}`
	if code := hook(t, w, body); code != http.StatusNoContent {
		t.Fatalf("hook = %d", code)
	}
	if code := hook(t, w, body); code != http.StatusNoContent {
		t.Fatalf("repeat hook = %d", code)
	}
	page, err := w.messages.Page(w.ctx, member, g.ID, chat.MessageQuery{Limit: 50})
	if err != nil {
		t.Fatalf("page: %v", err)
	}
	var posted []chat.Message
	for _, m := range page.Messages {
		if len(m.Files) > 0 {
			posted = append(posted, m)
		}
	}
	if len(posted) != 1 || posted[0].AuthorID == nil || *posted[0].AuthorID != owner.UserID || posted[0].Files[0].Size != 4096 || posted[0].Files[0].ContentType != "video/mp4" {
		t.Fatalf("posted = %+v", posted)
	}
}

func TestDirectCallsRecordOnlyWhenEveryoneMay(t *testing.T) {
	w := newWorld(t)
	a, b := w.user("a"), w.user("b")
	d, _, err := w.chats.Direct(w.ctx, a, chat.DirectInput{UserID: b.UserID})
	if err != nil {
		t.Fatalf("direct: %v", err)
	}
	if _, err := w.chats.Accept(w.ctx, b, d.ID); err != nil {
		t.Fatalf("accept: %v", err)
	}
	ticket, err := w.calls.Join(w.ctx, a, d.ID, chat.JoinInput{})
	if err != nil {
		t.Fatalf("join: %v", err)
	}
	w.settings.values[chat.RecordingSetting] = chat.RecordAdmins
	if _, err := w.records.Start(w.ctx, a, ticket.Call.ID, chat.RecordingInput{}); !errors.Is(err, chat.ErrRecordingOff) {
		t.Fatalf("admins rule in direct: %v", err)
	}
	w.settings.values[chat.RecordingSetting] = chat.RecordEveryone
	if call, err := w.records.Start(w.ctx, a, ticket.Call.ID, chat.RecordingInput{}); err != nil || !call.Recording {
		t.Fatalf("everyone = %+v, %v", call, err)
	}
	if _, err := w.records.Start(w.ctx, b, ticket.Call.ID, chat.RecordingInput{}); !errors.Is(err, chat.ErrCallNotFound) {
		t.Fatalf("not in call: %v", err)
	}
}

func TestRecordingFailsCleanlyWithoutEgress(t *testing.T) {
	w := newWorld(t)
	owner, member := w.user("owner"), w.user("member")
	g, _, err := w.chats.Create(w.ctx, owner, chat.ChatInput{Kind: chat.KindGroup, Title: "Команда", MemberIDs: []string{member.UserID, owner.UserID}})
	if err != nil {
		t.Fatalf("group: %v", err)
	}
	ticket, err := w.calls.Join(w.ctx, owner, g.ID, chat.JoinInput{})
	if err != nil {
		t.Fatalf("join: %v", err)
	}
	w.settings.values[chat.RecordingSetting] = chat.RecordEveryone
	w.egress.broken = true
	if _, err := w.records.Start(w.ctx, owner, ticket.Call.ID, chat.RecordingInput{}); !errors.Is(err, chat.ErrRecordingUnavailable) {
		t.Fatalf("broken egress: %v", err)
	}
	if call, err := w.calls.Active(w.ctx, owner, g.ID); err != nil || call.Recording {
		t.Fatalf("after failure = %+v, %v", call, err)
	}
	w.egress.broken = false
	if call, err := w.records.Start(w.ctx, owner, ticket.Call.ID, chat.RecordingInput{}); err != nil || !call.Recording {
		t.Fatalf("retry = %+v, %v", call, err)
	}
}
