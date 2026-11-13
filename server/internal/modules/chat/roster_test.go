package chat_test

import (
	"errors"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/modules/chat"
)

const courseID = "1c2d3e4f-5a6b-4c7d-8e9f-0a1b2c3d4e5f"

func TestCourseChatsFollowTheCourseList(t *testing.T) {
	w := newWorld(t)
	teacher, first, second := w.user("dana"), w.user("aigerim"), w.user("erlan")
	w.rosters.set(module.Roster{ID: courseID, Title: "Биология", Leads: []string{teacher.UserID}, Members: []string{first.UserID}})
	w.sync()
	w.sync()

	list, err := w.chats.List(w.ctx, first)
	if err != nil || len(list) != 1 || list[0].Kind != chat.KindChannel || list[0].Title != "Биология" || list[0].Role != chat.RoleMember {
		t.Fatalf("student chats = %+v, %v", list, err)
	}
	channel := list[0]
	if channel.Origin == nil || channel.Origin.Kind != "study.course" || channel.Origin.ID != courseID {
		t.Fatalf("origin = %+v", channel.Origin)
	}
	w.send(teacher, channel.ID, "t1", "Добро пожаловать на курс")
	if err := w.chats.RemoveMember(w.ctx, first, channel.ID, first.UserID); !errors.Is(err, chat.ErrSyncedMember) {
		t.Fatalf("synced leave: %v", err)
	}

	origin := &chat.Origin{Kind: "study.course", ID: courseID}
	if _, _, err := w.chats.Create(w.ctx, first, chat.ChatInput{Kind: chat.KindGroup, Origin: origin}); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("student creates course group: %v", err)
	}
	group, created, err := w.chats.Create(w.ctx, teacher, chat.ChatInput{Kind: chat.KindGroup, Origin: origin})
	if err != nil || !created || group.Title != "Биология" || group.MemberCount != 2 || group.Role != chat.RoleAdmin {
		t.Fatalf("course group = %+v, %v", group, err)
	}
	again, created, err := w.chats.Create(w.ctx, teacher, chat.ChatInput{Kind: chat.KindGroup, Origin: origin})
	if err != nil || created || again.ID != group.ID {
		t.Fatalf("second course group = %+v, %v", again, err)
	}
	w.send(first, group.ID, "s1", "Когда семинар?")
	if _, _, err := w.chats.Create(w.ctx, teacher, chat.ChatInput{Kind: chat.KindGroup, Origin: &chat.Origin{Kind: "teams.team", ID: courseID}}); !errors.Is(err, chat.ErrOriginNotFound) {
		t.Fatalf("unknown origin: %v", err)
	}

	w.rosters.set(module.Roster{ID: courseID, Title: "Биология", Leads: []string{teacher.UserID}, Members: []string{second.UserID}})
	w.sync()
	if _, err := w.chats.Get(w.ctx, first, channel.ID); !errors.Is(err, chat.ErrChatNotFound) {
		t.Fatalf("removed student still sees channel: %v", err)
	}
	mine, err := w.chats.List(w.ctx, second)
	if err != nil || len(mine) != 2 {
		t.Fatalf("new student chats = %+v, %v", mine, err)
	}

	w.rosters.set(module.Roster{ID: courseID, Title: "Биология", Leads: []string{teacher.UserID}, Members: []string{second.UserID}, Archived: true})
	w.sync()
	if _, _, err := w.messages.Send(w.ctx, teacher, channel.ID, chat.MessageInput{ClientID: "t2", Text: "Курс окончен"}); !errors.Is(err, chat.ErrReadOnly) {
		t.Fatalf("archived send: %v", err)
	}
	archived, err := w.chats.Get(w.ctx, second, group.ID)
	if err != nil || archived.ArchivedAt == nil {
		t.Fatalf("archived group = %+v, %v", archived, err)
	}
}

func TestGroupRostersBecomeGroupChats(t *testing.T) {
	w := newWorld(t)
	first, second := w.user("aigerim"), w.user("erlan")
	w.rosters.set(module.Roster{ID: courseID, Title: "Альфа", Leads: []string{}, Members: []string{first.UserID, second.UserID}, Group: true})
	w.sync()
	list, err := w.chats.List(w.ctx, second)
	if err != nil || len(list) != 1 || list[0].Kind != chat.KindGroup || list[0].Title != "Альфа" {
		t.Fatalf("team chats = %+v, %v", list, err)
	}
	w.send(second, list[0].ID, "s1", "Созвон в 18:00?")
	w.rosters.set(module.Roster{ID: courseID, Title: "Бета", Leads: []string{}, Members: []string{first.UserID, second.UserID}, Group: true})
	w.sync()
	if c, err := w.chats.Get(w.ctx, first, list[0].ID); err != nil || c.Title != "Бета" {
		t.Fatalf("renamed team chat = %+v, %v", c, err)
	}
}
