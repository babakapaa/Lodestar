package chat_test

import (
	"errors"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/modules/chat"
)

func TestHistoryPagesAndSyncsEditsAndDeletes(t *testing.T) {
	w := newWorld(t)
	a, b := w.user("aigerim"), w.user("erlan")
	d, _, err := w.chats.Direct(w.ctx, a, chat.DirectInput{UserID: b.UserID})
	if err != nil {
		t.Fatalf("direct: %v", err)
	}
	first := w.send(a, d.ID, "1", "первое")
	second := w.send(b, d.ID, "2", "второе")
	third := w.send(a, d.ID, "3", "третье")
	reply, _, err := w.messages.Send(w.ctx, b, d.ID, chat.MessageInput{ClientID: "4", Text: "ответ", ReplyTo: &first.Seq})
	if err != nil || *reply.ReplyTo != first.Seq {
		t.Fatalf("reply = %+v, %v", reply, err)
	}
	missing := int64(99)
	if _, _, err := w.messages.Send(w.ctx, b, d.ID, chat.MessageInput{ClientID: "5", Text: "x", ReplyTo: &missing}); !errors.Is(err, chat.ErrInvalidMessage) {
		t.Fatalf("reply to missing: %v", err)
	}

	before := reply.Seq
	page, err := w.messages.Page(w.ctx, a, d.ID, chat.MessageQuery{Before: &before, Limit: 2})
	if err != nil || !page.More || len(page.Messages) != 2 || page.Messages[0].Seq != second.Seq || page.Messages[1].Seq != third.Seq {
		t.Fatalf("page = %+v, %v", page, err)
	}

	if _, err := w.messages.Edit(w.ctx, a, second.ID, chat.EditInput{Text: "чужое"}); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("edit other: %v", err)
	}
	edited, err := w.messages.Edit(w.ctx, b, second.ID, chat.EditInput{Text: "второе, исправлено"})
	if err != nil || edited.EditedAt == nil || edited.Text != "второе, исправлено" || edited.Seq != second.Seq {
		t.Fatalf("edited = %+v, %v", edited, err)
	}
	if _, err := w.messages.Delete(w.ctx, b, third.ID); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("delete other in direct: %v", err)
	}
	gone, err := w.messages.Delete(w.ctx, a, third.ID)
	if err != nil || !gone.Deleted || gone.Text != "" {
		t.Fatalf("deleted = %+v, %v", gone, err)
	}

	since := reply.Rev
	changes, err := w.messages.Page(w.ctx, b, d.ID, chat.MessageQuery{Since: &since, Limit: 50})
	if err != nil || changes.More || len(changes.Messages) != 2 || changes.Messages[0].ID != second.ID || changes.Messages[1].ID != third.ID {
		t.Fatalf("changes = %+v, %v", changes, err)
	}
	if _, err := w.messages.Edit(w.ctx, a, third.ID, chat.EditInput{Text: "снова"}); !errors.Is(err, chat.ErrMessageNotFound) {
		t.Fatalf("edit deleted: %v", err)
	}
	list, err := w.chats.List(w.ctx, a)
	if err != nil || list[0].Seq != reply.Seq || list[0].Rev != gone.Rev {
		t.Fatalf("list = %+v, %v", list, err)
	}
}
