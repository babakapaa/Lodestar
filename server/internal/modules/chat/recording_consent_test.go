package chat_test

import (
	"errors"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/modules/chat"
)

func TestRecordingWaitsForEveryonesConsent(t *testing.T) {
	w := newWorld(t)
	w.settings.values[chat.RecordingSetting] = chat.RecordEveryone
	w.settings.values[chat.ConsentSetting] = true
	owner, member, late := w.user("owner"), w.user("member"), w.user("late")
	g, _, err := w.chats.Create(w.ctx, owner, chat.ChatInput{Kind: chat.KindGroup, Title: "Команда", MemberIDs: []string{member.UserID, late.UserID}})
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

	asked, err := w.records.Start(w.ctx, owner, callID, chat.RecordingInput{AudioOnly: true})
	if err != nil || asked.Recording || asked.Consent == nil || asked.Consent.RequestedBy != owner.UserID ||
		!asked.Consent.AudioOnly || len(asked.Consent.Waiting) != 1 || asked.Consent.Waiting[0] != member.UserID {
		t.Fatalf("asked = %+v, %v", asked, err)
	}
	if n := len(w.egress.started); n != 0 {
		t.Fatalf("recording started before consent: %d", n)
	}
	declined, err := w.records.Consent(w.ctx, member, callID, chat.ConsentInput{})
	if err != nil || declined.Consent != nil || declined.Recording {
		t.Fatalf("declined = %+v, %v", declined, err)
	}
	if _, err := w.records.Consent(w.ctx, member, callID, chat.ConsentInput{Agree: true}); !errors.Is(err, chat.ErrRecordingState) {
		t.Fatalf("consent without request: %v", err)
	}

	if _, err := w.records.Start(w.ctx, owner, callID, chat.RecordingInput{AudioOnly: true}); err != nil {
		t.Fatalf("ask again: %v", err)
	}
	started, err := w.records.Consent(w.ctx, member, callID, chat.ConsentInput{Agree: true})
	if err != nil || !started.Recording || !started.AudioOnly || started.Consent != nil {
		t.Fatalf("started = %+v, %v", started, err)
	}
	if got := w.egress.started[0]["audio_only"]; got != true {
		t.Fatalf("egress audio only = %v", got)
	}

	if _, err := w.calls.Join(w.ctx, late, g.ID, chat.JoinInput{}); !errors.Is(err, chat.ErrConsentRequired) {
		t.Fatalf("late joins without consent: %v", err)
	}
	joined, err := w.calls.Join(w.ctx, late, g.ID, chat.JoinInput{Consent: true})
	if err != nil || len(joined.Call.People) != 3 {
		t.Fatalf("late joins = %+v, %v", joined, err)
	}
}
