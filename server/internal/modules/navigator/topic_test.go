package navigator_test

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/modules/navigator"
)

func TestTopicsShowOnlyFromFiveStudents(t *testing.T) {
	w := newWorld(t)
	teacher := w.user("dana")
	w.teach(teacher)
	if _, err := w.materials.Add(w.ctx, teacher, courseID, navigator.MaterialInput{Title: "Строение клетки", Text: membrane}); err != nil {
		t.Fatalf("add: %v", err)
	}
	if err := navigator.NewIndexer(w.kernel).Index(w.ctx, w.tenant); err != nil {
		t.Fatalf("index: %v", err)
	}
	w.model.reply = "Мембрана «отделяет содержимое клетки от внешней среды» [1]. Как думаешь, зачем ей это?"
	ask := func(i int) {
		t.Helper()
		student := w.user(fmt.Sprintf("student-%d", i))
		w.enrol(module.ModeStudy, student)
		id := fmt.Sprintf("0b8a3c1e-5d2f-4a6b-9c7d-1e2f3a4b5c%02d", i)
		if _, _, err := w.navigator.Create(w.ctx, student, navigator.DialogInput{ID: id, ScopeKind: "assignment", ScopeID: assignmentID}); err != nil {
			t.Fatalf("create %d: %v", i, err)
		}
		for _, key := range []string{"q1", "q2"} {
			if _, _, err := w.navigator.Ask(w.ctx, student, id, navigator.MessageInput{Key: key, Text: "Что делает клеточная мембрана?", Intent: navigator.IntentUnderstand}, nil); err != nil {
				t.Fatalf("ask %d: %v", i, err)
			}
		}
	}
	for i := range 4 {
		ask(i)
	}
	few, err := w.topics.Course(w.ctx, teacher, courseID, 0)
	if err != nil || few.Students != nil || len(few.Topics) != 0 || few.Threshold != 5 || few.Weeks != 8 {
		t.Fatalf("four students = %+v, %v", few, err)
	}
	ask(4)
	ask(5)
	many, err := w.topics.Course(w.ctx, teacher, courseID, 4)
	if err != nil || many.Students == nil || *many.Students != 6 {
		t.Fatalf("six students = %+v, %v", many, err)
	}
	if len(many.Topics) != 1 || many.Topics[0].Title != "Строение клетки" || many.Topics[0].Students != 6 || many.Topics[0].Recent != 6 {
		t.Fatalf("topics = %+v", many.Topics)
	}
	if many.Suggestion == nil || many.Suggestion.MaterialID != many.Topics[0].MaterialID {
		t.Fatalf("suggestion = %+v", many.Suggestion)
	}
	if len(many.Intents) != 1 || many.Intents[0].Intent != navigator.IntentUnderstand || many.Intents[0].Students != 6 {
		t.Fatalf("intents = %+v", many.Intents)
	}
	last := many.Trend[len(many.Trend)-1]
	if len(many.Trend) != 4 || last.Students == nil || *last.Students != 6 || many.Trend[0].Students != nil {
		t.Fatalf("trend = %+v", many.Trend)
	}
	if _, err := w.topics.Course(w.ctx, w.user("student-0"), courseID, 0); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("student sees topics: %v", err)
	}
	if rec := w.get(teacher, "/v1/navigator/courses/"+courseID+"/topics?weeks=2"); rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"students":6`) {
		t.Fatalf("http = %d %s", rec.Code, rec.Body)
	}
	if rec := w.get(teacher, "/v1/navigator/courses/"+courseID+"/topics?weeks=x"); rec.Code != http.StatusBadRequest {
		t.Fatalf("bad weeks = %d", rec.Code)
	}
	members := append([]string{teacher.UserID}, w.course.students...)
	var metrics []module.Metric
	err = w.kernel.Tenants.Run(w.ctx, w.tenant, func(tx pgx.Tx) error {
		var err error
		metrics, err = navigator.NavigatorStats{}.Stats(w.ctx, tx, members, time.Now().Add(-time.Hour), time.Now().Add(time.Hour))
		return err
	})
	if err != nil || len(metrics) != 1 || metrics[0].Key != navigator.MetricUsers || metrics[0].Value != 6.0/7 {
		t.Fatalf("stats = %+v, %v", metrics, err)
	}
}
