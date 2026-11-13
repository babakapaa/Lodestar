package chat_test

import (
	"errors"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/modules/chat"
)

func TestChannelLetsOnlyAdminsWriteAndShowsViewsWithoutNames(t *testing.T) {
	w := newWorld(t)
	dean, student, other := w.user("dean"), w.user("aigerim"), w.user("erlan")
	in := chat.ChatInput{Kind: chat.KindChannel, Title: "Деканат", MemberIDs: []string{student.UserID, other.UserID}}

	if _, _, err := w.chats.Create(w.ctx, dean, in); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("create without permission: %v", err)
	}
	w.grant(dean, chat.CreateChannels)
	c, created, err := w.chats.Create(w.ctx, dean, in)
	if err != nil || !created || c.Role != chat.RoleOwner || c.MemberCount != 3 {
		t.Fatalf("channel = %+v, %v", c, err)
	}
	if _, _, err := w.messages.Send(w.ctx, student, c.ID, chat.MessageInput{ClientID: "s1", Text: "Можно вопрос?"}); !errors.Is(err, chat.ErrReadOnly) {
		t.Fatalf("subscriber send: %v", err)
	}
	m := w.send(dean, c.ID, "d1", "Сессия начинается 10 января")
	if _, err := w.chats.Read(w.ctx, student, c.ID, chat.ReadInput{Seq: m.Seq}); err != nil {
		t.Fatalf("read: %v", err)
	}

	page, err := w.messages.Page(w.ctx, dean, c.ID, chat.MessageQuery{Limit: 50})
	if err != nil || len(page.Messages) != 1 || page.Messages[0].Views == nil || *page.Messages[0].Views != 1 {
		t.Fatalf("page = %+v, %v", page, err)
	}
	seen, err := w.chats.Get(w.ctx, student, c.ID)
	if err != nil || len(seen.Members) != 1 || seen.Members[0].UserID != dean.UserID || seen.MemberCount != 3 {
		t.Fatalf("subscriber sees members = %+v, %v", seen, err)
	}
	full, err := w.chats.Get(w.ctx, dean, c.ID)
	if err != nil || len(full.Members) != 3 {
		t.Fatalf("owner sees members = %+v, %v", full, err)
	}
}

func TestGroupMembersMuteAndLeave(t *testing.T) {
	w := newWorld(t)
	owner, member, late := w.user("aigerim"), w.user("erlan"), w.user("dana")

	g, _, err := w.chats.Create(w.ctx, owner, chat.ChatInput{Kind: chat.KindGroup, Title: "Команда", MemberIDs: []string{member.UserID, owner.UserID}})
	if err != nil || g.MemberCount != 2 {
		t.Fatalf("group = %+v, %v", g, err)
	}
	if _, err := w.chats.AddMembers(w.ctx, member, g.ID, chat.MembersInput{UserIDs: []string{late.UserID}}); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("member adds: %v", err)
	}
	if _, err := w.chats.AddMembers(w.ctx, owner, g.ID, chat.MembersInput{UserIDs: []string{"00000000-0000-4000-8000-000000000000"}}); !errors.Is(err, access.ErrUserNotFound) {
		t.Fatalf("add unknown: %v", err)
	}
	w.send(owner, g.ID, "o1", "Старое сообщение")
	added, err := w.chats.AddMembers(w.ctx, owner, g.ID, chat.MembersInput{UserIDs: []string{late.UserID}})
	if err != nil || added.MemberCount != 3 {
		t.Fatalf("add = %+v, %v", added, err)
	}
	joined, err := w.chats.Get(w.ctx, late, g.ID)
	if err != nil || joined.Unread != 0 {
		t.Fatalf("late joiner = %+v, %v", joined, err)
	}
	if _, err := w.chats.Mute(w.ctx, late, g.ID, chat.MuteInput{Muted: true}); err != nil {
		t.Fatalf("mute: %v", err)
	}
	w.send(member, g.ID, "m1", "Встречаемся в 15:00")
	if n := w.sent.last(); len(n.UserIDs) != 1 || n.UserIDs[0] != owner.UserID || n.Link != "/chats/"+g.ID {
		t.Fatalf("notification = %+v", n)
	}

	if err := w.chats.RemoveMember(w.ctx, owner, g.ID, owner.UserID); !errors.Is(err, chat.ErrOwnerLeaves) {
		t.Fatalf("owner leaves: %v", err)
	}
	if err := w.chats.RemoveMember(w.ctx, member, g.ID, late.UserID); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("member removes other: %v", err)
	}
	if err := w.chats.RemoveMember(w.ctx, member, g.ID, member.UserID); err != nil {
		t.Fatalf("leave: %v", err)
	}
	if _, err := w.chats.Get(w.ctx, member, g.ID); !errors.Is(err, chat.ErrChatNotFound) {
		t.Fatalf("after leave: %v", err)
	}
	if w.events(member, chat.ChangedEvent) == 0 {
		t.Fatal("leaver was not told")
	}
}
