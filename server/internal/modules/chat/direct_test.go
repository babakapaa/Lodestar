package chat_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/modules/chat"
)

func TestDirectChatDeliversMessagesAndReadState(t *testing.T) {
	w := newWorld(t)
	aigerim, erlan, outsider := w.user("aigerim"), w.user("erlan"), w.user("dana")

	d, created, err := w.chats.Direct(w.ctx, aigerim, chat.DirectInput{UserID: erlan.UserID})
	if err != nil || !created || d.Kind != chat.KindDirect || d.Title != "erlan" || len(d.Members) != 2 {
		t.Fatalf("direct = %+v, %v, %v", d, created, err)
	}
	again, created, err := w.chats.Direct(w.ctx, erlan, chat.DirectInput{UserID: aigerim.UserID})
	if err != nil || created || again.ID != d.ID || again.Title != "aigerim" {
		t.Fatalf("again = %+v, %v, %v", again, created, err)
	}

	m := w.send(aigerim, d.ID, "c1", "Привет, ты сдал реферат?")
	if m.Seq != 1 || m.Text != "Привет, ты сдал реферат?" || *m.AuthorID != aigerim.UserID {
		t.Fatalf("message = %+v", m)
	}
	repeat, created, err := w.messages.Send(w.ctx, aigerim, d.ID, chat.MessageInput{ClientID: "c1", Text: "Привет, ты сдал реферат?"})
	if err != nil || created || repeat.ID != m.ID {
		t.Fatalf("repeat = %+v, %v, %v", repeat, created, err)
	}
	if bytes.Contains(w.raw(`SELECT data FROM chat.messages`), []byte("реферат")) {
		t.Fatal("message text is stored in the clear")
	}
	if w.events(erlan, chat.MessageEvent) != 1 {
		t.Fatal("erlan got no live event")
	}
	list, err := w.chats.List(w.ctx, erlan)
	if err != nil || len(list) != 1 || !list[0].Request || list[0].Unread != 1 || list[0].Last == nil || list[0].Last.Text != m.Text {
		t.Fatalf("list = %+v, %v", list, err)
	}
	if _, err := w.chats.Accept(w.ctx, erlan, d.ID); err != nil {
		t.Fatalf("accept: %v", err)
	}
	w.send(aigerim, d.ID, "c2", "Жду ответа")
	if n := w.sent.last(); n.Type != chat.NotifyMessage || len(n.UserIDs) != 1 || n.UserIDs[0] != erlan.UserID || n.Params["chat"] != "aigerim" {
		t.Fatalf("notification = %+v", n)
	}
	read, err := w.chats.Read(w.ctx, erlan, d.ID, chat.ReadInput{Seq: 99})
	if err != nil || read.Unread != 0 || read.ReadSeq != 2 {
		t.Fatalf("read = %+v, %v", read, err)
	}
	mine, err := w.chats.Get(w.ctx, aigerim, d.ID)
	if err != nil || mine.PeerReadSeq == nil || *mine.PeerReadSeq != 2 {
		t.Fatalf("peer read = %+v, %v", mine, err)
	}
	if w.events(aigerim, chat.ReadEvent) != 1 {
		t.Fatal("aigerim did not learn that erlan read the message")
	}

	if _, err := w.chats.Get(w.ctx, outsider, d.ID); !errors.Is(err, chat.ErrChatNotFound) {
		t.Fatalf("outsider get: %v", err)
	}
	if _, _, err := w.messages.Send(w.ctx, outsider, d.ID, chat.MessageInput{ClientID: "x", Text: "hi"}); !errors.Is(err, chat.ErrChatNotFound) {
		t.Fatalf("outsider send: %v", err)
	}
	if _, _, err := w.chats.Direct(w.ctx, aigerim, chat.DirectInput{UserID: aigerim.UserID}); !errors.Is(err, chat.ErrInvalidChat) {
		t.Fatalf("self direct: %v", err)
	}
}
