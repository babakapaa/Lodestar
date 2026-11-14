package chat_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/modules/chat"
)

func recordCall(t *testing.T, w *world) (string, string) {
	t.Helper()
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
	if _, err := w.records.Start(w.ctx, owner, ticket.Call.ID, chat.RecordingInput{AudioOnly: true}); err != nil {
		t.Fatalf("start: %v", err)
	}
	if _, err := w.records.Stop(w.ctx, owner, ticket.Call.ID); err != nil {
		t.Fatalf("stop: %v", err)
	}
	body := `{"event":"egress_ended","egressInfo":{"egressId":"EG_1","roomName":"` + ticket.Room + `","status":"EGRESS_COMPLETE","fileResults":[{"filename":"` + w.egress.key(0) + `","size":"2048"}]}}`
	if code := hook(t, w, body); code != http.StatusNoContent {
		t.Fatalf("hook = %d", code)
	}
	return g.ID, member.UserID
}

func (w *world) attached(chatID, userID string) []string {
	w.t.Helper()
	page, err := w.messages.Page(w.ctx, auth.Principal{TenantID: w.tenant, UserID: userID}, chatID, chat.MessageQuery{Limit: 50})
	if err != nil {
		w.t.Fatalf("page: %v", err)
	}
	var names []string
	for _, m := range page.Messages {
		for _, f := range m.Files {
			names = append(names, f.Name+" "+f.ContentType)
		}
	}
	return names
}

func TestTranscriptIsPostedNextToRecording(t *testing.T) {
	w := newWorld(t)
	w.settings.values[chat.TranscriptSetting] = true
	w.speech.text = " Добрый день, начинаем защиту. "
	chatID, member := recordCall(t, w)

	if err := w.scribe.RunTenant(w.ctx, w.tenant); err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(w.speech.heard) != 1 || !strings.HasSuffix(w.speech.heard[0], ".ogg:audio") {
		t.Fatalf("heard = %v", w.speech.heard)
	}
	names := w.attached(chatID, member)
	if len(names) != 2 || !strings.Contains(strings.Join(names, ","), "-transcript.txt text/plain") {
		t.Fatalf("files = %v", names)
	}
	var stored []string
	for key, v := range w.speech.objects {
		if strings.HasPrefix(key, w.tenant+"/") {
			stored = append(stored, v)
		}
	}
	if len(stored) != 1 || stored[0] != "Добрый день, начинаем защиту." {
		t.Fatalf("stored = %v", stored)
	}
	if err := w.scribe.RunTenant(w.ctx, w.tenant); err != nil || len(w.speech.heard) != 1 {
		t.Fatalf("second run heard %v, %v", w.speech.heard, err)
	}
}

func TestTranscriptsFollowUniversityRule(t *testing.T) {
	w := newWorld(t)
	chatID, member := recordCall(t, w)
	if err := w.scribe.RunTenant(w.ctx, w.tenant); err != nil || len(w.speech.heard) != 0 {
		t.Fatalf("off: heard %v, %v", w.speech.heard, err)
	}
	if names := w.attached(chatID, member); len(names) != 1 {
		t.Fatalf("files = %v", names)
	}
}

func TestBrokenSpeechServerRetriesLater(t *testing.T) {
	w := newWorld(t)
	w.settings.values[chat.TranscriptSetting] = true
	w.speech.broken = true
	chatID, member := recordCall(t, w)
	if err := w.scribe.RunTenant(w.ctx, w.tenant); err != nil {
		t.Fatalf("run: %v", err)
	}
	if err := w.scribe.RunTenant(w.ctx, w.tenant); err != nil || len(w.speech.heard) != 1 {
		t.Fatalf("retried at once: %v, %v", w.speech.heard, err)
	}
	if names := w.attached(chatID, member); len(names) != 1 {
		t.Fatalf("files = %v", names)
	}
}
