package chat_test

import (
	"errors"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/modules/chat"
)

func TestChannelCallIsStartedByAdmins(t *testing.T) {
	w := newWorld(t)
	teacher, student := w.user("teacher"), w.user("aigerim")
	w.grant(teacher, chat.CreateChannels)
	c, _, err := w.chats.Create(w.ctx, teacher, chat.ChatInput{Kind: chat.KindChannel, Title: "Биология", MemberIDs: []string{student.UserID}})
	if err != nil {
		t.Fatalf("channel: %v", err)
	}
	if _, err := w.calls.Join(w.ctx, student, c.ID); !errors.Is(err, chat.ErrCallNotStarted) {
		t.Fatalf("student starts: %v", err)
	}
	lead, err := w.calls.Join(w.ctx, teacher, c.ID)
	if err != nil || lead.Muted {
		t.Fatalf("teacher starts = %+v, %v", lead, err)
	}
	if n := w.sent.last(); n.Type != chat.NotifyCall || len(n.UserIDs) != 1 || n.UserIDs[0] != student.UserID {
		t.Fatalf("ring = %+v", n)
	}
	joined, err := w.calls.Join(w.ctx, student, c.ID)
	if err != nil || !joined.Muted || joined.Call.ID != lead.Call.ID || len(joined.Call.People) != 2 {
		t.Fatalf("student joins = %+v, %v", joined, err)
	}
}
