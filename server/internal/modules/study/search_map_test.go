package study_test

import (
	"errors"
	"slices"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/modules/study"
)

func TestSearchMapIsDraftedPrivatelyAndChecked(t *testing.T) {
	w := newWorld(t)
	group := w.unit(nil, "group", "ИС-21")
	teacher := w.user("dana", "dana@kaznu.kz")
	student := w.user("aigerim", "aigerim@kaznu.kz")
	w.join(group, student)
	c := w.course(nil, "Биология")
	w.add(study.ListUnits, c, group.ID)
	w.add(study.ListTeachers, c, teacher.UserID)

	if _, err := w.work.Create(w.ctx, teacher, c.ID, study.AssignmentInput{
		Title: "Клетка", Grading: study.GradingPass, SearchMap: &study.MapRules{MinScholarly: 101},
	}); !errors.Is(err, study.ErrInvalidAssignment) {
		t.Fatalf("bad rules = %v", err)
	}
	a, err := w.work.Create(w.ctx, teacher, c.ID, study.AssignmentInput{
		Title: "Клетка", Grading: study.GradingPass, Published: true,
		SearchMap: &study.MapRules{MinSources: 2, MinScholarly: 50},
	})
	if err != nil || a.SearchMap == nil || a.SearchMap.MinSources != 2 {
		t.Fatalf("create = %+v, %v", a, err)
	}

	empty, err := w.work.SearchMap(w.ctx, student, a.ID)
	if err != nil || empty.ChangedAt != nil || empty.Rules.MinScholarly != 50 {
		t.Fatalf("empty draft = %+v, %v", empty, err)
	}
	if _, err := w.work.SearchMap(w.ctx, teacher, a.ID); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("teacher reads draft = %v", err)
	}

	web := study.MapSource{Reference: "Сайт о клетке", Kind: "web", Found: "поиск", Why: "Обзор"}
	m := study.SearchMap{Used: []study.MapSource{web}, Queries: []string{" строение клетки ", ""}}
	d, err := w.work.SaveSearchMap(w.ctx, student, a.ID, m)
	if err != nil || d.ChangedAt == nil || !slices.Equal(d.Map.Queries, []string{"строение клетки"}) {
		t.Fatalf("save = %+v, %v", d, err)
	}
	if !slices.Equal(d.Missing, []string{study.MapMissingSources, study.MapMissingScholarly, study.MapMissingPath}) {
		t.Fatalf("missing = %v", d.Missing)
	}

	in := study.SubmissionInput{ClientID: "6f1d7a52-3c38-4a8e-9a43-1f0e5f3c1a02", Body: "Ответ", SearchMap: &d.Map}
	if _, _, err := w.work.Submit(w.ctx, student, a.ID, in); !errors.Is(err, study.ErrSearchMapIncomplete) {
		t.Fatalf("incomplete submit = %v", err)
	}

	article := study.MapSource{Reference: "Иванов. Строение клетки", Kind: "article", DOI: "10.1000/xyz", Why: "Свежие данные", Taken: "Схема"}
	m = study.SearchMap{
		Used:     []study.MapSource{web, article},
		Rejected: []study.MapSource{{Reference: "Старый учебник", Kind: "book", Why: "Устарел"}},
		Path:     "Сначала искал про ядро, потом понял, что важнее мембрана.",
	}
	if d, err = w.work.SaveSearchMap(w.ctx, student, a.ID, m); err != nil || len(d.Missing) != 0 {
		t.Fatalf("complete = %+v, %v", d, err)
	}
	in.SearchMap = &d.Map
	sub, _, err := w.work.Submit(w.ctx, student, a.ID, in)
	if err != nil || sub.SearchMap == nil || len(sub.SearchMap.Rejected) != 1 {
		t.Fatalf("submit = %+v, %v", sub, err)
	}
	work, err := w.work.Work(w.ctx, teacher, a.ID)
	if err != nil || work[0].Latest == nil || work[0].Latest.SearchMap.Path != m.Path {
		t.Fatalf("teacher sees = %+v, %v", work, err)
	}

	long := study.SearchMap{Used: []study.MapSource{{Reference: ""}}}
	if _, err := w.work.SaveSearchMap(w.ctx, student, a.ID, long); !errors.Is(err, study.ErrInvalidSearchMap) {
		t.Fatalf("source without reference = %v", err)
	}
}
