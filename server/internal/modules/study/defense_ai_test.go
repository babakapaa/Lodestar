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

func TestDefenseQuestionsAndProtocolAreDraftedFromTheWork(t *testing.T) {
	d := newDefenseCase(t, study.DefenseFull)
	w := d.w
	if _, err := d.defenses.DraftQuestions(w.ctx, d.teacher, d.work.ID, d.first.UserID); !errors.Is(err, study.ErrAIUnavailable) {
		t.Fatalf("no model: %v", err)
	}

	model := &scriptedModel{}
	k := w.kernel()
	k.Model = model
	defenses := study.NewDefenseService(k)
	if _, err := defenses.DraftQuestions(w.ctx, d.teacher, d.work.ID, d.first.UserID); !errors.Is(err, study.ErrNoWork) {
		t.Fatalf("no work: %v", err)
	}

	m := study.SearchMap{
		Used:     []study.MapSource{{Reference: "Codd E. F. A relational model of data", Kind: "article", Found: "eLibrary", Why: "Основа теории"}},
		Rejected: []study.MapSource{{Reference: "Блог про SQL", Kind: "web", Found: "Google", Why: "Нет ссылок на источники"}},
		Path:     "Начала с учебника, потом нашла первоисточник.",
	}
	if _, _, err := w.work.Submit(w.ctx, d.first, d.work.ID, study.SubmissionInput{
		ClientID: "4c2b1f0e-7a8d-4e3b-9c51-2d6f0a1b3e77", Body: "Третья нормальная форма убирает транзитивные зависимости.", SearchMap: &m,
	}); err != nil {
		t.Fatalf("submit: %v", err)
	}

	model.reply = "```json\n" + `{"questions":[
		{"kind":"path","text":"Почему вы отбросили блог про SQL?"},
		{"kind":"understanding","text":"Что такое транзитивная зависимость?","mark":3,"note":"лишнее"},
		{"kind":"guess","text":"Списали ли вы работу?"},
		{"kind":"analysis","text":"Почему 3НФ, а не НФБК?"},
		{"kind":"application","text":"Как нормализовать таблицу заказов?"},
		{"kind":"analysis","text":"Чем первоисточник лучше учебника?"},
		{"kind":"understanding","text":"Что такое ключ?"},
		{"kind":"understanding","text":"Что такое 2НФ?"},
		{"kind":"path","text":"Как вы нашли статью Кодда?"}
	]}` + "\n```"
	questions, err := defenses.DraftQuestions(w.ctx, d.teacher, d.work.ID, d.first.UserID)
	if err != nil || len(questions) != 7 {
		t.Fatalf("questions = %+v, %v", questions, err)
	}
	for _, q := range questions {
		if q.Kind == "guess" || q.Mark != nil || q.Note != "" {
			t.Fatalf("draft kept %+v", q)
		}
	}
	sent := model.lastUser()
	for _, want := range []string{"транзитивные зависимости", "Codd E. F.", "Нет ссылок на источники", "первоисточник"} {
		if !strings.Contains(sent, want) {
			t.Fatalf("model did not get %q:\n%s", want, sent)
		}
	}
	if _, err := defenses.DraftQuestions(w.ctx, d.second, d.work.ID, d.first.UserID); err == nil {
		t.Fatal("student drafted questions")
	}

	model.reply = "  Студентка уверенно объяснила нормальные формы.  "
	mark := 3
	summary, err := defenses.DraftProtocol(w.ctx, d.teacher, d.work.ID, d.first.UserID, study.DefenseInput{
		Questions: []study.DefenseQuestion{{Kind: study.QuestionUnderstanding, Text: "Что такое 3НФ?", Mark: &mark, Note: "Без подсказок"}},
	})
	if err != nil || summary != "Студентка уверенно объяснила нормальные формы." {
		t.Fatalf("protocol = %q, %v", summary, err)
	}
	if sent := model.lastUser(); !strings.Contains(sent, "Отметка: 3") || !strings.Contains(sent, "Без подсказок") {
		t.Fatalf("protocol input:\n%s", sent)
	}

	model.reply = "не JSON"
	if _, err := defenses.DraftQuestions(w.ctx, d.teacher, d.work.ID, d.first.UserID); !errors.Is(err, study.ErrAIFailed) {
		t.Fatalf("bad answer: %v", err)
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
