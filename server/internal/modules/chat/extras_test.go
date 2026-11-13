package chat_test

import (
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/search"
	"github.com/babakapaa/lodestar/server/internal/modules/chat"
)

func TestFilesTravelWithMessagesAndOnlyMembersRead(t *testing.T) {
	w := newWorld(t)
	a, b, outsider := w.user("aigerim"), w.user("erlan"), w.user("dana")
	g, _, err := w.chats.Create(w.ctx, a, chat.ChatInput{Kind: chat.KindGroup, Title: "Проект", MemberIDs: []string{b.UserID}})
	if err != nil {
		t.Fatalf("group: %v", err)
	}
	fileID := w.file(a, "план.txt")
	foreign := w.file(b, "чужой.txt")

	if _, _, err := w.messages.Send(w.ctx, a, g.ID, chat.MessageInput{ClientID: "x", FileIDs: []string{foreign}}); !errors.Is(err, chat.ErrBadFile) {
		t.Fatalf("foreign file: %v", err)
	}
	m, _, err := w.messages.Send(w.ctx, a, g.ID, chat.MessageInput{ClientID: "f1", FileIDs: []string{fileID}})
	if err != nil || len(m.Files) != 1 || m.Files[0].Name != "план.txt" || m.Text != "" {
		t.Fatalf("file message = %+v, %v", m, err)
	}
	page, err := w.messages.Page(w.ctx, b, g.ID, chat.MessageQuery{Limit: 50})
	if err != nil || len(page.Messages) != 1 || len(page.Messages[0].Files) != 1 {
		t.Fatalf("page = %+v, %v", page, err)
	}
	link := chat.FileLink{}
	err = w.tx.Run(w.ctx, w.tenant, func(tx pgx.Tx) error {
		if ok, err := link.Readable(w.ctx, tx, b, fileID); err != nil || !ok {
			t.Errorf("member cannot read: %v %v", ok, err)
		}
		if ok, err := link.Readable(w.ctx, tx, outsider, fileID); err != nil || ok {
			t.Errorf("outsider reads: %v %v", ok, err)
		}
		if ok, err := link.InUse(w.ctx, tx, fileID); err != nil || !ok {
			t.Errorf("not in use: %v %v", ok, err)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("link: %v", err)
	}
	if _, err := w.messages.Delete(w.ctx, a, m.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	_ = w.tx.Run(w.ctx, w.tenant, func(tx pgx.Tx) error {
		if ok, _ := link.InUse(w.ctx, tx, fileID); ok {
			t.Error("deleted message still holds the file")
		}
		return nil
	})
}

func TestMentionsRingEvenWhenMuted(t *testing.T) {
	w := newWorld(t)
	a, b, c := w.user("aigerim"), w.user("erlan"), w.user("dana")
	g, _, err := w.chats.Create(w.ctx, a, chat.ChatInput{Kind: chat.KindGroup, Title: "ИС-21", MemberIDs: []string{b.UserID, c.UserID}})
	if err != nil {
		t.Fatalf("group: %v", err)
	}
	if _, err := w.chats.Mute(w.ctx, b, g.ID, chat.MuteInput{Muted: true}); err != nil {
		t.Fatalf("mute: %v", err)
	}
	outsider := w.user("outsider")
	m, _, err := w.messages.Send(w.ctx, a, g.ID, chat.MessageInput{ClientID: "1", Text: "@erlan глянь", Mentions: []string{b.UserID, outsider.UserID}})
	if err != nil || len(m.Mentions) != 1 || m.Mentions[0] != b.UserID {
		t.Fatalf("message = %+v, %v", m, err)
	}
	var mention, plain []string
	for _, n := range w.sent.list {
		switch n.Type {
		case chat.NotifyMention:
			mention = n.UserIDs
		case chat.NotifyMessage:
			plain = n.UserIDs
		}
	}
	if len(mention) != 1 || mention[0] != b.UserID || len(plain) != 1 || plain[0] != c.UserID {
		t.Fatalf("mention = %v, plain = %v", mention, plain)
	}
	edited, err := w.messages.Edit(w.ctx, a, m.ID, chat.EditInput{Text: "@erlan посмотри"})
	if err != nil || len(edited.Mentions) != 1 {
		t.Fatalf("edited = %+v, %v", edited, err)
	}
}

func TestReportCarriesOnlyChosenMessages(t *testing.T) {
	w := newWorld(t)
	a, b, moderator := w.user("aigerim"), w.user("erlan"), w.user("moderator")
	d := w.direct(b, a)
	w.send(b, d.ID, "1", "Обычное сообщение")
	bad := w.send(b, d.ID, "2", "Оскорбление")

	if _, err := w.reports.Create(w.ctx, w.user("dana"), chat.ReportInput{ChatID: d.ID, MessageIDs: []string{bad.ID}, Reason: "x"}); !errors.Is(err, chat.ErrChatNotFound) {
		t.Fatalf("outsider report: %v", err)
	}
	r, err := w.reports.Create(w.ctx, a, chat.ReportInput{ChatID: d.ID, MessageIDs: []string{bad.ID}, Reason: "Оскорбляет"})
	if err != nil || len(r.Messages) != 1 || r.Messages[0].Text != "Оскорбление" || r.Messages[0].AuthorName != "erlan" {
		t.Fatalf("report = %+v, %v", r, err)
	}
	if _, err := w.reports.Open(w.ctx, a); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("student reads reports: %v", err)
	}
	w.grant(moderator, chat.ReviewReports)
	open, err := w.reports.Open(w.ctx, moderator)
	if err != nil || len(open) != 1 {
		t.Fatalf("open = %+v, %v", open, err)
	}
	if _, err := w.reports.Resolve(w.ctx, moderator, r.ID); err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if open, _ := w.reports.Open(w.ctx, moderator); len(open) != 0 {
		t.Fatalf("still open: %+v", open)
	}
}

func TestSearchFindsMyGroupsAndChannels(t *testing.T) {
	w := newWorld(t)
	a, b := w.user("aigerim"), w.user("erlan")
	if _, _, err := w.chats.Create(w.ctx, a, chat.ChatInput{Kind: chat.KindGroup, Title: "Олимпиада по биологии"}); err != nil {
		t.Fatalf("group: %v", err)
	}
	q, _ := search.ParseQuery("биологии")
	for _, tc := range []struct {
		who  string
		want int
	}{{"owner", 1}, {"other", 0}} {
		p := a
		if tc.who == "other" {
			p = b
		}
		err := w.tx.Run(w.ctx, w.tenant, func(tx pgx.Tx) error {
			hits, err := chat.ChatsSource{}.Search(w.ctx, tx, p, q, 10)
			if err != nil || len(hits) != tc.want {
				t.Errorf("%s hits = %+v, %v", tc.who, hits, err)
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestOldMessagesExpireByUniversityRule(t *testing.T) {
	w := newWorld(t)
	a, b := w.user("aigerim"), w.user("erlan")
	d := w.direct(a, b)
	old := w.send(a, d.ID, "1", "Давнее")
	fresh := w.send(a, d.ID, "2", "Свежее")
	if _, err := w.pool.Exec(w.ctx, `UPDATE chat.messages SET created_at = $2 WHERE id = $1`, old.ID, time.Now().AddDate(0, 0, -40)); err != nil {
		t.Fatalf("age: %v", err)
	}
	if err := w.purger.Purge(w.ctx, w.tenant); err != nil {
		t.Fatalf("purge off: %v", err)
	}
	w.settings.retention = 30
	if err := w.purger.Purge(w.ctx, w.tenant); err != nil {
		t.Fatalf("purge: %v", err)
	}
	since := fresh.Rev
	changes, err := w.messages.Page(w.ctx, b, d.ID, chat.MessageQuery{Since: &since, Limit: 50})
	if err != nil || len(changes.Messages) != 1 || changes.Messages[0].ID != old.ID || !changes.Messages[0].Deleted {
		t.Fatalf("changes = %+v, %v", changes, err)
	}
}
