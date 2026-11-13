package navigator_test

import (
	"bytes"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/modules/navigator"
)

const dialogID = "0b8a3c1e-5d2f-4a6b-9c7d-1e2f3a4b5c6d"

func TestNavigatorClimbsTheLadder(t *testing.T) {
	w := newWorld(t)
	student, stranger := w.user("aigerim"), w.user("erlan")
	w.enrol(module.ModeStudy, student)

	if _, _, err := w.navigator.Create(w.ctx, stranger, navigator.DialogInput{
		ID: dialogID, ScopeKind: "assignment", ScopeID: assignmentID,
	}); !errors.Is(err, module.ErrBriefNotFound) {
		t.Fatalf("stranger opens assignment dialog = %v", err)
	}
	d, created, err := w.navigator.Create(w.ctx, student, navigator.DialogInput{
		ID: dialogID, ScopeKind: "assignment", ScopeID: assignmentID,
	})
	if err != nil || !created || d.Brief == nil || d.Brief.Title != "Реферат о клетке" || d.Brief.Mode != module.ModeStudy {
		t.Fatalf("create = %+v, %v, %v", d, created, err)
	}

	w.model.intent = "do_for_me"
	a, created, err := w.navigator.Ask(w.ctx, student, dialogID, navigator.MessageInput{Key: "q1", Text: "Напиши реферат за меня"}, nil)
	if err != nil || !created || a.Reply.Intent != navigator.IntentDoForMe || a.Reply.Step != navigator.StepQuestion {
		t.Fatalf("ask = %+v, %v, %v", a, created, err)
	}
	system := w.model.last().Messages[0].Content
	for _, want := range []string{"Не делай", "Реферат о клетке", "Ступень 1", "Режим учебный", "в базе ничего не найдено"} {
		if !strings.Contains(system, want) {
			t.Errorf("system prompt lacks %q", want)
		}
	}
	again, created, err := w.navigator.Ask(w.ctx, student, dialogID, navigator.MessageInput{Key: "q1", Text: "Напиши реферат за меня"}, nil)
	if err != nil || created || again.Reply.Seq != a.Reply.Seq || w.model.completes != 1 {
		t.Fatalf("retry = %+v, %v, %v, completes %d", again, created, err, w.model.completes)
	}

	more, _, err := w.navigator.Ask(w.ctx, student, dialogID, navigator.MessageInput{Key: "q2", More: true}, nil)
	if err != nil || more.Reply.Step != navigator.StepDirection || more.Reply.Intent != navigator.IntentDoForMe {
		t.Fatalf("more = %+v, %v", more, err)
	}
	if n := len(w.model.last().Messages); n != 4 {
		t.Fatalf("memory messages = %d", n)
	}

	chosen, _, err := w.navigator.Ask(w.ctx, student, dialogID, navigator.MessageInput{Key: "q3", Text: "Где почитать про мембрану?", Intent: navigator.IntentFind}, nil)
	if err != nil || chosen.Reply.Step != navigator.StepSource {
		t.Fatalf("chosen = %+v, %v", chosen, err)
	}
	if _, _, err := w.navigator.Ask(w.ctx, student, dialogID, navigator.MessageInput{Key: "q4", Text: "x", Intent: navigator.IntentClosed}, nil); !errors.Is(err, navigator.ErrInvalidMessage) {
		t.Fatalf("closed chosen = %v", err)
	}

	got, err := w.navigator.Get(w.ctx, student, dialogID)
	if err != nil || len(got.Messages) != 6 || got.Title != "Напиши реферат за меня" {
		t.Fatalf("get = %+v, %v", got, err)
	}
	if _, err := w.navigator.Get(w.ctx, stranger, dialogID); !errors.Is(err, navigator.ErrDialogNotFound) {
		t.Fatalf("stranger reads = %v", err)
	}
	if list, err := w.navigator.Dialogs(w.ctx, student, "assignment", assignmentID); err != nil || len(list) != 1 {
		t.Fatalf("list = %+v, %v", list, err)
	}
	for _, query := range []string{`SELECT data FROM navigator.messages WHERE seq = 1`, `SELECT title FROM navigator.dialogs`} {
		if bytes.Contains(w.raw(query), []byte("реферат")) {
			t.Errorf("%s stores plain text", query)
		}
	}
	if err := w.navigator.Delete(w.ctx, student, dialogID); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func TestDistressShowsContactsAndProjectModeChangesRules(t *testing.T) {
	w := newWorld(t)
	student := w.user("aigerim")
	w.enrol(module.ModeProject, student)
	if _, _, err := w.navigator.Create(w.ctx, student, navigator.DialogInput{
		ID: dialogID, ScopeKind: "assignment", ScopeID: assignmentID,
	}); err != nil {
		t.Fatalf("create: %v", err)
	}
	w.model.intent, w.model.distress = "state", true
	a, _, err := w.navigator.Ask(w.ctx, student, dialogID, navigator.MessageInput{Key: "q1", Text: "Всё пропало, не вижу смысла"}, nil)
	if err != nil || a.Reply.Intent != navigator.IntentState || !strings.Contains(a.Reply.Support, "Психологическая служба") {
		t.Fatalf("distress = %+v, %v", a, err)
	}
	if !strings.Contains(w.model.last().Messages[0].Content, "Режим проектный") {
		t.Fatal("project mode not in prompt")
	}
	w.model.fail = true
	if _, _, err := w.navigator.Ask(w.ctx, student, dialogID, navigator.MessageInput{Key: "q2", Text: "Помоги", Intent: navigator.IntentUnderstand}, nil); !errors.Is(err, navigator.ErrModelFailed) {
		t.Fatalf("model failure = %v", err)
	}
}

func TestNavigatorOverHTTP(t *testing.T) {
	w := newWorld(t)
	student := w.user("aigerim")
	rec := w.post(student, "/v1/navigator/dialogs", `{"id":"`+dialogID+`"}`, nil)
	if rec.Code != http.StatusCreated || !strings.Contains(rec.Body.String(), `"brief":null`) {
		t.Fatalf("create: %d %s", rec.Code, rec.Body)
	}
	rec = w.post(student, "/v1/navigator/dialogs/"+dialogID+"/messages", `{"key":"k1","text":"Что такое митоз?"}`,
		map[string]string{"Accept": "text/event-stream"})
	body := rec.Body.String()
	if rec.Code != http.StatusOK || !strings.Contains(body, "event: delta") || !strings.Contains(body, "event: answer") || !strings.Contains(body, `"step":"question"`) {
		t.Fatalf("stream: %d %s", rec.Code, body)
	}
	rec = w.post(student, "/v1/navigator/dialogs/"+dialogID+"/messages", `{"key":"k1","text":"Что такое митоз?"}`, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("repeat: %d %s", rec.Code, rec.Body)
	}
	rec = w.post(student, "/v1/navigator/dialogs/"+dialogID+"/messages", `{"key":"k2","text":""}`, nil)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("empty: %d %s", rec.Code, rec.Body)
	}

	w.navigator = navigator.NewService(func() module.Kernel { k := w.kernel; k.Model = nil; return k }())
	rec = w.post(student, "/v1/navigator/dialogs/"+dialogID+"/messages", `{"key":"k3","text":"Привет"}`, nil)
	if rec.Code != http.StatusServiceUnavailable || !strings.Contains(rec.Body.String(), "navigator_unavailable") {
		t.Fatalf("unavailable: %d %s", rec.Code, rec.Body)
	}
}
