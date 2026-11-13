package attendance_test

import (
	"errors"
	"net/netip"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/modules/attendance"
)

func TestTeacherKeepsJournalAndSeesLongAbsence(t *testing.T) {
	w := newWorld(t)
	start := time.Date(2026, 9, 8, 9, 0, 0, 0, time.UTC)
	var ids []string
	for i := range 3 {
		day := start.AddDate(0, 0, 7*i)
		l := w.lesson(day.Format("2006-01-02"), day)
		w.clock.now = day.Add(5 * time.Minute)
		s, err := w.service.Open(w.ctx, w.teacher, l.SeriesID, l.Date)
		if err != nil {
			t.Fatalf("open %d: %v", i, err)
		}
		ids = append(ids, s.ID)
		if i < 2 {
			if _, err := w.service.Check(w.ctx, w.aigerim, attendance.CheckInput{Code: s.Code}, netip.Addr{}); err != nil {
				t.Fatalf("check %d: %v", i, err)
			}
		}
		closed, err := w.service.Close(w.ctx, w.teacher, s.ID)
		if err != nil || closed.Open || closed.Counts.Unmarked != 0 {
			t.Fatalf("close %d = %+v, %v", i, closed, err)
		}
	}
	cancelled := w.lesson("2026-09-29", start.AddDate(0, 0, 21))
	w.lessons.list[len(w.lessons.list)-1].Cancelled = true
	if _, err := w.service.Open(w.ctx, w.teacher, cancelled.SeriesID, cancelled.Date); !errors.Is(err, attendance.ErrLessonCancelled) {
		t.Fatalf("cancelled: %v", err)
	}
	if _, err := w.service.Close(w.ctx, w.teacher, ids[0]); !errors.Is(err, attendance.ErrSessionClosed) {
		t.Fatalf("close twice: %v", err)
	}
	rec := w.call(w.teacher, "PUT", "/v1/attendance/sessions/"+ids[2]+"/marks/"+w.aigerim.UserID, `{"status":"excused"}`)
	if rec.Code != 200 || !strings.Contains(rec.Body.String(), `"manual"`) {
		t.Fatalf("correct = %d %s", rec.Code, rec.Body)
	}
	if rec := w.call(w.teacher, "PUT", "/v1/attendance/sessions/"+ids[2]+"/marks/"+w.stranger.UserID, `{"status":"present"}`); rec.Code != 403 {
		t.Fatalf("stranger mark = %d %s", rec.Code, rec.Body)
	}
	if rec := w.call(w.teacher, "PUT", "/v1/attendance/sessions/"+ids[2]+"/marks/"+w.aigerim.UserID, `{"status":"sleeping"}`); rec.Code != 400 {
		t.Fatalf("bad status = %d %s", rec.Code, rec.Body)
	}

	if _, err := w.service.Course(w.ctx, w.aigerim, "course-1", "2026-09-01", "2026-10-01"); !errors.Is(err, attendance.ErrLessonNotFound) {
		t.Fatalf("student sees journal: %v", err)
	}
	c, err := w.service.Course(w.ctx, w.teacher, "course-1", "2026-09-01", "2026-10-01")
	if err != nil || len(c.Lessons) != 4 || c.Lessons[3].SessionID != nil || len(c.Students) != 2 {
		t.Fatalf("course = %+v, %v", c, err)
	}
	aigerim, erlan := c.Students[0], c.Students[1]
	if aigerim.Counts.Present != 2 || aigerim.Counts.Excused != 1 || aigerim.LongAbsent || *aigerim.Rate != 1 {
		t.Fatalf("aigerim = %+v", aigerim)
	}
	if erlan.Counts.Absent != 3 || erlan.Streak != 3 || !erlan.LongAbsent || *erlan.Rate != 0 {
		t.Fatalf("erlan = %+v", erlan)
	}
	if _, err := w.service.Course(w.ctx, w.admin, "course-1", "2026-10-01", "2026-09-01"); !errors.Is(err, attendance.ErrInvalidRange) {
		t.Fatalf("range: %v", err)
	}
	if rec := w.call(w.teacher, "GET", "/v1/attendance/courses/course-1/export?format=csv&from=2026-09-01&to=2026-10-01", ""); rec.Code != 200 ||
		!strings.Contains(rec.Body.String(), "Ерлан,absent,absent,absent,0,0,3,0,0") {
		t.Fatalf("export = %d %s", rec.Code, rec.Body)
	}
	mine, err := w.service.Mine(w.ctx, w.aigerim)
	if err != nil || len(mine) != 3 || mine[0].Status != attendance.Excused || mine[2].Method != attendance.ByCode {
		t.Fatalf("mine = %+v, %v", mine, err)
	}
	var metrics []module.Metric
	err = w.tx.Run(w.ctx, w.tenant, func(tx pgx.Tx) error {
		var err error
		metrics, err = attendance.AttendanceStats{}.Stats(w.ctx, tx, []string{w.aigerim.UserID, w.erlan.UserID}, start, start.AddDate(0, 1, 0))
		return err
	})
	if err != nil || len(metrics) != 1 || metrics[0].Key != attendance.MetricRate || metrics[0].Value != 0.4 {
		t.Fatalf("stats = %+v, %v", metrics, err)
	}
}
