package study_test

import (
	"errors"
	"slices"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/modules/study"
)

func TestCourseRosterTakesStudentsFromGroupsAndList(t *testing.T) {
	w := newWorld(t)
	stream := w.unit(nil, "stream", "ИС-2023")
	group := w.unit(&stream, "group", "ИС-21")
	fromGroup := w.user("aigerim", "aigerim@kaznu.kz")
	listed := w.user("erlan", "erlan@kaznu.kz")
	teacher := w.user("dana", "dana@kaznu.kz")
	w.join(group, fromGroup)
	w.join(group, teacher)
	c := w.course(nil, "Базы данных")
	w.add(study.ListUnits, c, stream.ID)
	w.add(study.ListStudents, c, listed.UserID)
	w.add(study.ListTeachers, c, teacher.UserID)
	empty := w.course(nil, "Пустой курс")

	source, ok := module.Rosters(study.Module{}.RosterSources(w.kernel())).Find(study.RosterCourse)
	if !ok {
		t.Fatal("course roster source is missing")
	}
	var all []module.Roster
	var one module.Roster
	var missing error
	err := w.tx.Run(w.ctx, w.tenant, func(tx pgx.Tx) error {
		var err error
		if all, err = source.All(w.ctx, tx); err != nil {
			return err
		}
		if one, err = source.One(w.ctx, tx, c.ID); err != nil {
			return err
		}
		_, missing = source.One(w.ctx, tx, "00000000-0000-4000-8000-000000000000")
		return nil
	})
	if err != nil {
		t.Fatalf("rosters: %v", err)
	}
	if len(all) != 2 || !errors.Is(missing, module.ErrRosterNotFound) {
		t.Fatalf("all = %+v, missing = %v", all, missing)
	}
	want := []string{fromGroup.UserID, listed.UserID}
	slices.Sort(want)
	if one.Title != "Базы данных" || !slices.Equal(one.Members, want) || !slices.Equal(one.Leads, []string{teacher.UserID}) {
		t.Fatalf("roster = %+v", one)
	}
	for _, r := range all {
		if r.ID == empty.ID && (len(r.Members) != 0 || len(r.Leads) != 0) {
			t.Fatalf("empty roster = %+v", r)
		}
	}
}
