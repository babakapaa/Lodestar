package study_test

import (
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/modules/study"
)

func TestPilotIncludesCourseMembers(t *testing.T) {
	w := newWorld(t)
	group := w.unit(nil, "group", "ИС-21")
	fromGroup := w.user("aigerim", "aigerim@kaznu.kz")
	teacher := w.user("dana", "dana@kaznu.kz")
	stranger := w.user("erlan", "erlan@kaznu.kz")
	w.join(group, fromGroup)
	c := w.course(nil, "Базы данных")
	other := w.course(nil, "Химия")
	w.add(study.ListUnits, c, group.ID)
	w.add(study.ListTeachers, c, teacher.UserID)

	source := study.Module{}.PilotSources(w.kernel())[0]
	err := w.tx.Run(w.ctx, w.tenant, func(tx pgx.Tx) error {
		courses, err := source.PilotCourses(w.ctx, tx, w.tenant)
		if err != nil {
			return err
		}
		if len(courses) != 2 {
			t.Fatalf("courses = %+v", courses)
		}
		for _, c := range []struct {
			user    string
			courses []string
			want    bool
		}{
			{fromGroup.UserID, []string{c.ID}, true},
			{teacher.UserID, []string{c.ID}, true},
			{stranger.UserID, []string{c.ID}, false},
			{teacher.UserID, []string{other.ID}, false},
		} {
			in, err := source.InPilot(w.ctx, tx, c.user, c.courses)
			if err != nil {
				return err
			}
			if in != c.want {
				t.Fatalf("in pilot %s %v = %v", c.user, c.courses, in)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("pilot: %v", err)
	}
}
