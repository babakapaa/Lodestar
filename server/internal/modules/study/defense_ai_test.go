package study_test

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/babakapaa/lodestar/server/internal/modules/study"
	"github.com/babakapaa/lodestar/server/internal/platform/llm"
)

type scriptedModel struct {
	mu       sync.Mutex
	reply    string
	requests []llm.Request
}

func (m *scriptedModel) Complete(_ context.Context, req llm.Request) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.requests = append(m.requests, req)
	return m.reply, nil
}

func (m *scriptedModel) Stream(ctx context.Context, req llm.Request, delta func(string) error) (string, error) {
	full, err := m.Complete(ctx, req)
	if err != nil {
		return "", err
	}
	return full, delta(full)
}

func (m *scriptedModel) lastUser() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	msgs := m.requests[len(m.requests)-1].Messages
	return msgs[len(msgs)-1].Content
}

func TestDefenseProtocolIsDraftedFromTheTeachersMarks(t *testing.T) {
	d := newDefenseCase(t, study.DefenseFull)
	w := d.w
	in := study.DefenseInput{
		Questions: []study.DefenseQuestion{{Kind: study.QuestionUnderstanding, Text: "Что такое 3НФ?", Mark: new(int), Note: "Без подсказок"}},
	}
	*in.Questions[0].Mark = 3
	if _, err := d.defenses.DraftProtocol(w.ctx, d.teacher, d.work.ID, d.first.UserID, in); !errors.Is(err, study.ErrAIUnavailable) {
		t.Fatalf("no model: %v", err)
	}

	model := &scriptedModel{}
	k := w.kernel()
	k.Model = model
	defenses := study.NewDefenseService(k)

	model.reply = "  Студентка уверенно объяснила нормальные формы.  "
	summary, err := defenses.DraftProtocol(w.ctx, d.teacher, d.work.ID, d.first.UserID, in)
	if err != nil || summary != "Студентка уверенно объяснила нормальные формы." {
		t.Fatalf("protocol = %q, %v", summary, err)
	}
	if sent := model.lastUser(); !strings.Contains(sent, "Отметка: 3") || !strings.Contains(sent, "Без подсказок") {
		t.Fatalf("protocol input:\n%s", sent)
	}
	if _, err := defenses.DraftProtocol(w.ctx, d.second, d.work.ID, d.first.UserID, in); err == nil {
		t.Fatal("student drafted a protocol")
	}

	model.reply = "   "
	if _, err := defenses.DraftProtocol(w.ctx, d.teacher, d.work.ID, d.first.UserID, in); !errors.Is(err, study.ErrAIFailed) {
		t.Fatalf("empty answer: %v", err)
	}
}

func TestStudentsAreRemindedBeforeTheirDefenseOnce(t *testing.T) {
	d := newDefenseCase(t, study.DefenseShort)
	w := d.w
	soon, err := d.defenses.AddSlots(w.ctx, d.teacher, d.work.ID, study.SlotsInput{StartsAt: time.Now().Add(time.Hour), Minutes: 10, Count: 2})
	if err != nil {
		t.Fatalf("slots: %v", err)
	}
	later, err := d.defenses.AddSlots(w.ctx, d.teacher, d.work.ID, study.SlotsInput{StartsAt: time.Now().Add(48 * time.Hour), Minutes: 10})
	if err != nil {
		t.Fatalf("later slot: %v", err)
	}
	if _, err := d.defenses.Book(w.ctx, d.first, d.work.ID, study.BookingInput{SlotID: soon[0].ID}); err != nil {
		t.Fatalf("book soon: %v", err)
	}
	if _, err := d.defenses.Book(w.ctx, d.second, d.work.ID, study.BookingInput{SlotID: later[0].ID}); err != nil {
		t.Fatalf("book later: %v", err)
	}
	reminder := study.NewDefenseReminder(w.kernel())
	for range 2 {
		if err := reminder.RemindAll(w.ctx); err != nil {
			t.Fatalf("remind: %v", err)
		}
	}
	got := w.sent.of(study.NotifyDefenseSoon)
	if len(got) != 1 || got[0].UserIDs[0] != d.first.UserID {
		t.Fatalf("reminders = %+v", got)
	}

	if _, err := d.defenses.Book(w.ctx, d.first, d.work.ID, study.BookingInput{SlotID: soon[1].ID}); err != nil {
		t.Fatalf("move: %v", err)
	}
	if err := reminder.RemindAll(w.ctx); err != nil {
		t.Fatalf("remind after move: %v", err)
	}
	if n := len(w.sent.of(study.NotifyDefenseSoon)); n != 2 {
		t.Fatalf("reminders after move = %d", n)
	}
}
