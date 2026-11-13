package setup_test

import (
	"errors"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/core/setup"
	"github.com/babakapaa/lodestar/server/internal/core/structure"
)

func (w *world) ready() {
	w.t.Helper()
	if err := w.setup.Begin(w.ctx, w.tenant); err != nil {
		w.t.Fatalf("begin: %v", err)
	}
	for _, q := range setup.LegalQuestions {
		if _, err := w.setup.Answer(w.ctx, w.admin, q.Key, setup.AnswerInput{Answer: "Решено", Basis: "Приказ"}); err != nil {
			w.t.Fatalf("answer %s: %v", q.Key, err)
		}
		if _, err := w.setup.ApproveAnswer(w.ctx, w.admin, q.Key); err != nil {
			w.t.Fatalf("approve %s: %v", q.Key, err)
		}
	}
	if _, err := w.setup.Sign(w.ctx, w.admin); err != nil {
		w.t.Fatalf("sign: %v", err)
	}
	for _, key := range []string{setup.SectionDeploy, setup.SectionUniversity, setup.SectionPeople, setup.SectionAdmins, setup.SectionLegal} {
		if _, err := w.setup.Review(w.ctx, w.admin, key, setup.ReviewInput{Approve: true}); err != nil {
			w.t.Fatalf("approve %s: %v", key, err)
		}
	}
	if _, err := w.structure.Create(w.ctx, w.admin, structure.UnitInput{Kind: "faculty", Name: "Биофак"}); err != nil {
		w.t.Fatalf("unit: %v", err)
	}
}

func TestPilotOpensOnlyChosenCourses(t *testing.T) {
	w := newWorld(t)
	w.ready()
	inside, outside := w.user(), w.user()
	bio, chem := "11111111-1111-1111-1111-111111111111", "22222222-2222-2222-2222-222222222222"
	w.pilot.courses = []module.PilotCourse{{ID: bio, Title: "Биология"}, {ID: chem, Title: "Химия"}}
	w.pilot.members[bio] = []string{inside.UserID}

	courses, err := w.setup.PilotCourses(w.ctx, w.admin)
	if err != nil || len(courses) != 2 {
		t.Fatalf("courses = %+v, %v", courses, err)
	}
	if _, err := w.setup.PilotCourses(w.ctx, inside); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("courses by student: %v", err)
	}
	for _, bad := range [][]string{{"33333333-3333-3333-3333-333333333333"}, {bio, bio}} {
		if _, err := w.setup.Launch(w.ctx, w.admin, setup.LaunchInput{Pilot: bad}); !errors.Is(err, setup.ErrInvalidPilot) {
			t.Fatalf("launch %v: %v", bad, err)
		}
	}
	o, err := w.setup.Launch(w.ctx, w.admin, setup.LaunchInput{Pilot: []string{bio}})
	if err != nil || !o.Launched || len(o.Pilot) != 1 || o.Pilot[0] != bio {
		t.Fatalf("launch = %+v, %v", o, err)
	}
	people := []struct {
		name string
		p    auth.Principal
		ok   bool
	}{{"inside", inside, true}, {"outside", outside, false}, {"admin", w.admin, true}}
	for _, c := range people {
		allowed, err := w.setup.Allowed(w.ctx, c.p)
		if err != nil || allowed != c.ok {
			t.Fatalf("%s allowed = %v, %v", c.name, allowed, err)
		}
		view, err := w.setup.Status(w.ctx, c.p)
		if err != nil || !view.Launched || !view.Pilot || view.Included != c.ok {
			t.Fatalf("%s status = %+v, %v", c.name, view, err)
		}
	}
	if _, err := w.setup.Expand(w.ctx, inside); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("expand by student: %v", err)
	}
	if o, err := w.setup.Expand(w.ctx, w.admin); err != nil || len(o.Pilot) != 0 {
		t.Fatalf("expand = %+v, %v", o, err)
	}
	if allowed, err := w.setup.Allowed(w.ctx, outside); err != nil || !allowed {
		t.Fatalf("outside after expand = %v, %v", allowed, err)
	}
	if _, err := w.setup.Expand(w.ctx, w.admin); !errors.Is(err, setup.ErrNotPilot) {
		t.Fatalf("expand twice: %v", err)
	}
}
