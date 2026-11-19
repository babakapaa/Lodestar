package navigator_test

import (
	"bytes"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/modules/navigator"
)

func TestHelpSummaryCountsStepsAndSources(t *testing.T) {
	w := newWorld(t)
	teacher, student := w.user("dana"), w.user("aigerim")
	w.teach(teacher)
	w.enrol(module.ModeStudy, student)
	if _, err := w.materials.Add(w.ctx, teacher, courseID, navigator.MaterialInput{Title: "Строение клетки", Reference: "Иванов И. И.", Text: membrane}); err != nil {
		t.Fatalf("add: %v", err)
	}
	if err := navigator.NewIndexer(w.kernel).Index(w.ctx, w.tenant); err != nil {
		t.Fatalf("index: %v", err)
	}
	w.model.reply = "Мембрана «отделяет содержимое клетки от внешней среды» [1]. Что дальше?"
	if _, _, err := w.navigator.Create(w.ctx, student, navigator.DialogInput{ID: dialogID, ScopeKind: "assignment", ScopeID: assignmentID}); err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, _, err := w.navigator.Ask(w.ctx, student, dialogID, navigator.MessageInput{Key: "q1", Text: "Что делает клеточная мембрана?", Intent: navigator.IntentFind}, nil); err != nil {
		t.Fatalf("ask: %v", err)
	}
	if _, _, err := w.navigator.Ask(w.ctx, student, dialogID, navigator.MessageInput{Key: "q2", More: true}, nil); err != nil {
		t.Fatalf("more: %v", err)
	}

	sum, err := w.navigator.Help(w.ctx, student, "assignment", assignmentID)
	if err != nil || sum.Replies != 2 || sum.Steps.Source != 1 || sum.Steps.Fragment != 1 || len(sum.Sources) != 1 || sum.Sources[0] != "Иванов И. И." {
		t.Fatalf("help = %+v, %v", sum, err)
	}
	if err := w.navigator.Delete(w.ctx, student, dialogID); err != nil {
		t.Fatalf("delete dialog: %v", err)
	}
	if kept, err := w.navigator.Help(w.ctx, student, "assignment", assignmentID); err != nil || kept.Replies != 2 {
		t.Fatalf("help after deleting the dialog = %+v, %v", kept, err)
	}
	if bytes.Contains(w.raw(`SELECT data FROM navigator.help`), []byte("Иванов")) {
		t.Fatal("help stores plain text")
	}
	if none, err := w.navigator.Help(w.ctx, teacher, "assignment", assignmentID); err != nil || none.Replies != 0 || none.Sources == nil {
		t.Fatalf("teacher help = %+v, %v", none, err)
	}
	rec := w.get(student, "/v1/navigator/help?scopeKind=assignment&scopeId="+assignmentID)
	if rec.Code != 200 || !bytes.Contains(rec.Body.Bytes(), []byte(`"replies":2`)) {
		t.Fatalf("help over HTTP = %d %s", rec.Code, rec.Body)
	}
}
