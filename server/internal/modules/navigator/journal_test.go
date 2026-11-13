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

func TestJournalKeepsQuestionsAndSources(t *testing.T) {
	w := newWorld(t)
	teacher, student, quiet := w.user("dana"), w.user("aigerim"), w.user("erlan")
	w.teach(teacher)
	w.enrol(module.ModeStudy, student)
	w.enrol(module.ModeStudy, quiet)
	w.prefs.off = []string{quiet.UserID}
	if _, err := w.materials.Add(w.ctx, teacher, courseID, navigator.MaterialInput{Title: "Строение клетки", Reference: "Иванов И. И.", Text: membrane}); err != nil {
		t.Fatalf("add: %v", err)
	}
	if err := navigator.NewIndexer(w.kernel).Index(w.ctx, w.tenant); err != nil {
		t.Fatalf("index: %v", err)
	}
	w.model.reply = "Мембрана «отделяет содержимое клетки от внешней среды» [1]. Что дальше?"
	for _, p := range []struct {
		name string
		id   string
	}{{"student", dialogID}, {"quiet", "7c1d2e3f-4a5b-4c6d-8e7f-9a0b1c2d3e4f"}} {
		who := student
		if p.name == "quiet" {
			who = quiet
		}
		if _, _, err := w.navigator.Create(w.ctx, who, navigator.DialogInput{ID: p.id, ScopeKind: "assignment", ScopeID: assignmentID}); err != nil {
			t.Fatalf("create: %v", err)
		}
		if _, _, err := w.navigator.Ask(w.ctx, who, p.id, navigator.MessageInput{Key: "q1", Text: "Что делает клеточная мембрана?", Intent: navigator.IntentFind}, nil); err != nil {
			t.Fatalf("ask: %v", err)
		}
		if _, _, err := w.navigator.Ask(w.ctx, who, p.id, navigator.MessageInput{Key: "q2", More: true}, nil); err != nil {
			t.Fatalf("more: %v", err)
		}
	}

	list, err := w.navigator.Journal(w.ctx, student, "assignment", assignmentID)
	if err != nil || len(list) != 1 || list[0].Query != "Что делает клеточная мембрана?" || len(list[0].Sources) != 1 || list[0].Sources[0].Reference != "Иванов И. И." {
		t.Fatalf("journal = %+v, %v", list, err)
	}
	if bytes.Contains(w.raw(`SELECT data FROM navigator.journal`), []byte("мембрана")) {
		t.Fatal("journal stores plain text")
	}
	if off, err := w.navigator.Journal(w.ctx, quiet, "", ""); err != nil || len(off) != 0 {
		t.Fatalf("journal while off = %+v, %v", off, err)
	}
	if err := w.navigator.Forget(w.ctx, quiet, list[0].ID); !errors.Is(err, navigator.ErrJournalNotFound) {
		t.Fatalf("forget other's entry = %v", err)
	}
	if err := w.navigator.Forget(w.ctx, student, list[0].ID); err != nil {
		t.Fatalf("forget: %v", err)
	}
}

func TestClosedAssignmentAndSupport(t *testing.T) {
	w := newWorld(t)
	student := w.user("aigerim")
	w.enrol(module.ModeClosed, student)
	if _, _, err := w.navigator.Create(w.ctx, student, navigator.DialogInput{ID: dialogID, ScopeKind: "assignment", ScopeID: assignmentID}); err != nil {
		t.Fatalf("create: %v", err)
	}
	w.navigator = navigator.NewService(func() module.Kernel { k := w.kernel; k.Model = nil; return k }())
	a, _, err := w.navigator.Ask(w.ctx, student, dialogID, navigator.MessageInput{Key: "q1", Text: "Реши задачу 3"}, nil)
	if err != nil || a.Reply.Intent != navigator.IntentClosed || !strings.Contains(a.Reply.Text, "контрольной") {
		t.Fatalf("closed = %+v, %v", a, err)
	}
	if list, _ := w.navigator.Journal(w.ctx, student, "", ""); len(list) != 0 {
		t.Fatalf("closed question journaled: %+v", list)
	}

	rec := w.get(student, "/v1/navigator/support")
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "Психологическая служба") {
		t.Fatalf("support: %d %s", rec.Code, rec.Body)
	}
}
