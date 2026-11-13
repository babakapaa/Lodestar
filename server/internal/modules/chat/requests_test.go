package chat_test

import (
	"errors"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/modules/chat"
)

func TestStrangerLandsInRequestsAndCanBeBlocked(t *testing.T) {
	w := newWorld(t)
	stranger, student, classmate := w.user("stranger"), w.user("aigerim"), w.user("erlan")

	d := w.direct(stranger, student)
	if d.Request {
		t.Fatal("the writer sees a request")
	}
	w.send(stranger, d.ID, "s1", "Привет, продашь конспект?")
	if n := w.sent.last(); n.Type == chat.NotifyMessage {
		t.Fatalf("request rang: %+v", n)
	}
	inbox, err := w.chats.List(w.ctx, student)
	if err != nil || len(inbox) != 1 || !inbox[0].Request {
		t.Fatalf("inbox = %+v, %v", inbox, err)
	}
	if err := w.blocks.Block(w.ctx, student, chat.DirectInput{UserID: stranger.UserID}); err != nil {
		t.Fatalf("block: %v", err)
	}
	if _, _, err := w.messages.Send(w.ctx, stranger, d.ID, chat.MessageInput{ClientID: "s2", Text: "Ответь"}); !errors.Is(err, chat.ErrBlocked) {
		t.Fatalf("blocked send: %v", err)
	}
	blocked, err := w.blocks.List(w.ctx, student)
	if err != nil || len(blocked) != 1 || blocked[0].UserID != stranger.UserID {
		t.Fatalf("blocked = %+v, %v", blocked, err)
	}
	if err := w.blocks.Unblock(w.ctx, student, stranger.UserID); err != nil {
		t.Fatalf("unblock: %v", err)
	}
	accepted, err := w.chats.Accept(w.ctx, student, d.ID)
	if err != nil || accepted.Request {
		t.Fatalf("accept = %+v, %v", accepted, err)
	}
	w.send(stranger, d.ID, "s3", "Спасибо")
	if n := w.sent.last(); n.Type != chat.NotifyMessage || n.UserIDs[0] != student.UserID {
		t.Fatalf("after accept = %+v", n)
	}

	if _, _, err := w.chats.Create(w.ctx, student, chat.ChatInput{Kind: chat.KindGroup, Title: "ИС-21", MemberIDs: []string{classmate.UserID}}); err != nil {
		t.Fatalf("group: %v", err)
	}
	known := w.direct(classmate, student)
	mine, err := w.chats.Get(w.ctx, student, known.ID)
	if err != nil || mine.Request {
		t.Fatalf("classmate chat is a request: %+v, %v", mine, err)
	}
}

func TestReplyAcceptsRequest(t *testing.T) {
	w := newWorld(t)
	a, b := w.user("aigerim"), w.user("erlan")
	d := w.direct(a, b)
	w.send(a, d.ID, "1", "Здравствуйте")
	w.send(b, d.ID, "2", "Добрый день")
	got, err := w.chats.Get(w.ctx, b, d.ID)
	if err != nil || got.Request {
		t.Fatalf("after reply = %+v, %v", got, err)
	}
}
