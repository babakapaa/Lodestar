package navigator_test

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/modules/navigator"
)

const membrane = `Клеточная мембрана отделяет содержимое клетки от внешней среды и регулирует обмен веществ.

Митохондрии вырабатывают энергию для клетки в форме АТФ.`

func TestTeacherMaterialsAreIndexed(t *testing.T) {
	w := newWorld(t)
	teacher, student, stranger := w.user("dana"), w.user("aigerim"), w.user("erlan")
	w.teach(teacher)
	w.enrol(module.ModeStudy, student)

	in := navigator.MaterialInput{Title: "Строение клетки", Reference: "Иванов И. И. Биология. 2020.", Text: membrane}
	if _, err := w.materials.Add(w.ctx, student, courseID, in); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("student adds = %v", err)
	}
	if _, err := w.materials.Add(w.ctx, teacher, courseID, navigator.MaterialInput{Title: "x", Text: "y", URL: "ftp://x"}); !errors.Is(err, navigator.ErrInvalidMaterial) {
		t.Fatalf("bad link = %v", err)
	}
	m, err := w.materials.Add(w.ctx, teacher, courseID, in)
	if err != nil || m.Passages != 1 || m.Indexed != 0 || m.CourseID != courseID {
		t.Fatalf("add = %+v, %v", m, err)
	}
	if err := navigator.NewIndexer(w.kernel).Index(w.ctx, w.tenant); err != nil {
		t.Fatalf("index: %v", err)
	}
	list, err := w.materials.List(w.ctx, student, courseID)
	if err != nil || len(list) != 1 || list[0].Indexed != 1 {
		t.Fatalf("list = %+v, %v", list, err)
	}
	if _, err := w.materials.List(w.ctx, stranger, courseID); !errors.Is(err, module.ErrBriefNotFound) {
		t.Fatalf("stranger lists = %v", err)
	}
	if err := w.materials.Delete(w.ctx, student, m.ID); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("student deletes = %v", err)
	}
	if err := w.materials.Delete(w.ctx, teacher, m.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if err := w.materials.Delete(w.ctx, teacher, m.ID); !errors.Is(err, navigator.ErrMaterialNotFound) {
		t.Fatalf("delete twice = %v", err)
	}

	long := strings.Repeat("Абзац про ядро клетки и его оболочку. ", 25)
	big, err := w.materials.Add(w.ctx, teacher, courseID, navigator.MaterialInput{Title: "Ядро", Text: long + "\n\n" + long + "\n\n" + long})
	if err != nil || big.Passages != 3 {
		t.Fatalf("big = %+v, %v", big, err)
	}

	rec := w.post(teacher, "/v1/navigator/courses/"+courseID+"/materials", `{"title":"Ткани","text":"Эпителий покрывает тело."}`, nil)
	if rec.Code != http.StatusCreated || !strings.Contains(rec.Body.String(), `"passages":1`) {
		t.Fatalf("http add: %d %s", rec.Code, rec.Body)
	}
}

func TestAnswersStayWithinFragments(t *testing.T) {
	w := newWorld(t)
	teacher, student := w.user("dana"), w.user("aigerim")
	w.teach(teacher)
	w.enrol(module.ModeStudy, student)
	if _, err := w.materials.Add(w.ctx, teacher, courseID, navigator.MaterialInput{Title: "Строение клетки", Text: membrane}); err != nil {
		t.Fatalf("add: %v", err)
	}
	if err := navigator.NewIndexer(w.kernel).Index(w.ctx, w.tenant); err != nil {
		t.Fatalf("index: %v", err)
	}
	if _, _, err := w.navigator.Create(w.ctx, student, navigator.DialogInput{ID: dialogID, ScopeKind: "assignment", ScopeID: assignmentID}); err != nil {
		t.Fatalf("create: %v", err)
	}
	ask := func(key, reply string) navigator.Message {
		t.Helper()
		w.model.reply = reply
		a, _, err := w.navigator.Ask(w.ctx, student, dialogID, navigator.MessageInput{Key: key, Text: "Что делает клеточная мембрана?", Intent: navigator.IntentUnderstand}, nil)
		if err != nil {
			t.Fatalf("ask %s: %v", key, err)
		}
		return a.Reply
	}

	var stages []navigator.Stage
	w.model.reply = "Мембрана «отделяет содержимое клетки от внешней среды» [1]. Как думаешь, зачем ей это?"
	a, _, err := w.navigator.Ask(w.ctx, student, dialogID, navigator.MessageInput{Key: "q1", Text: "Что делает клеточная мембрана?", Intent: navigator.IntentUnderstand},
		func(s navigator.Stage) error { stages = append(stages, s); return nil })
	if err != nil {
		t.Fatalf("ask: %v", err)
	}
	if len(a.Reply.Citations) != 1 || a.Reply.Citations[0].Title != "Строение клетки" || a.Reply.Text != w.model.reply {
		t.Fatalf("grounded = %+v", a.Reply)
	}
	if strings.Join([]string{string(stages[0]), string(stages[1]), string(stages[2])}, ",") != "search,write,check" {
		t.Fatalf("stages = %v", stages)
	}
	if system := w.model.last().Messages[0].Content; !strings.Contains(system, "[1] Строение клетки") || !strings.Contains(system, "Митохондрии") {
		t.Fatalf("prompt lacks fragment: %s", system)
	}

	invented := ask("q2", "Мембрана «состоит из кварков и глюонов» [1]. Что ты знаешь о белках?")
	if invented.Text != "Что ты знаешь о белках?" || len(invented.Citations) != 0 {
		t.Fatalf("invented quote kept: %+v", invented)
	}
	foreign := ask("q3", "Это описано в учебнике [7]. Что ещё нашёл?")
	if foreign.Text != "Что ещё нашёл?" {
		t.Fatalf("unknown citation kept: %+v", foreign)
	}

	w.model.unsupported = []int{1}
	flagged := ask("q4", "Мембрана состоит из целлюлозы [1]. Проверь по фрагменту.")
	if flagged.Text != "Проверь по фрагменту." {
		t.Fatalf("unsupported kept: %+v", flagged)
	}
	w.model.unsupported = nil

	w.model.leak = true
	before := w.model.completes
	leaked := ask("q5", "Вот готовый реферат целиком.")
	if w.model.completes-before != 2 || strings.Contains(leaked.Text, "реферат") || !strings.Contains(w.model.last().Messages[0].Content, "не прошёл проверку") {
		t.Fatalf("leak = %+v, completes %d", leaked, w.model.completes-before)
	}

	w.enrol(module.ModeProject, student)
	whole := ask("q6", "Вот готовый код задания целиком.")
	if strings.Contains(whole.Text, "готовый код") || !strings.Contains(w.model.verifier().Messages[0].Content, "Режим проектный") {
		t.Fatalf("project mode leak check = %+v", whole)
	}
	w.model.leak = false
	project := ask("q7", "Вот пример кода для разбора.")
	if project.Text != "Вот пример кода для разбора." {
		t.Fatalf("project mode snippet = %+v", project)
	}
}
