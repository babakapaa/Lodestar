package attendance_test

import (
	"errors"
	"net/netip"
	"strings"
	"testing"
	"time"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/modules/attendance"
)

func TestStudentsCheckInByQRAndCode(t *testing.T) {
	w := newWorld(t)
	l := w.lesson("2026-09-08", time.Date(2026, 9, 8, 9, 0, 0, 0, time.UTC))
	if _, err := w.service.Open(w.ctx, w.stranger, l.SeriesID, l.Date); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("stranger opens: %v", err)
	}
	if _, err := w.service.Open(w.ctx, w.teacher, l.SeriesID, "2026-09-09"); !errors.Is(err, attendance.ErrLessonNotFound) {
		t.Fatalf("no lesson: %v", err)
	}
	s, err := w.service.Open(w.ctx, w.teacher, l.SeriesID, l.Date)
	if err != nil || !s.Open || len(s.Code) != 4 || s.Counts.Unmarked != 2 || len(s.Marks) != 2 {
		t.Fatalf("open = %+v, %v", s, err)
	}
	if _, err := w.service.Session(w.ctx, w.aigerim, s.ID); !errors.Is(err, attendance.ErrSessionNotFound) {
		t.Fatalf("student sees session: %v", err)
	}
	open, err := w.service.OpenLessons(w.ctx, w.aigerim)
	if err != nil || len(open) != 1 || open[0].Status != nil || !open[0].Methods.QR {
		t.Fatalf("open lessons = %+v, %v", open, err)
	}

	qr, err := w.service.QR(w.ctx, w.teacher, s.ID)
	if err != nil || qr.Token == "" {
		t.Fatalf("qr = %+v, %v", qr, err)
	}
	w.clock.move(16 * time.Second)
	res, err := w.service.Check(w.ctx, w.aigerim, attendance.CheckInput{SessionID: s.ID, Token: qr.Token}, netip.Addr{})
	if err != nil || res.Status != attendance.Present || res.CourseTitle != "Сети" {
		t.Fatalf("qr check = %+v, %v", res, err)
	}
	w.clock.move(30 * time.Second)
	if _, err := w.service.Check(w.ctx, w.erlan, attendance.CheckInput{SessionID: s.ID, Token: qr.Token}, netip.Addr{}); !errors.Is(err, attendance.ErrWrongCode) {
		t.Fatalf("old qr: %v", err)
	}
	w.clock.move(2 * time.Minute)
	if _, err := w.service.Check(w.ctx, w.erlan, attendance.CheckInput{Code: s.Code}, netip.Addr{}); !errors.Is(err, attendance.ErrWrongCode) {
		t.Fatalf("expired code: %v", err)
	}
	w.clock.move(8 * time.Minute)
	fresh, err := w.service.RefreshCode(w.ctx, w.teacher, s.ID)
	if err != nil {
		t.Fatalf("refresh: %v", err)
	}
	if _, err := w.service.Check(w.ctx, w.stranger, attendance.CheckInput{Code: fresh.Code}, netip.Addr{}); !errors.Is(err, attendance.ErrWrongCode) {
		t.Fatalf("stranger code: %v", err)
	}
	late, err := w.service.Check(w.ctx, w.erlan, attendance.CheckInput{Code: fresh.Code}, netip.Addr{})
	if err != nil || late.Status != attendance.Late {
		t.Fatalf("late = %+v, %v", late, err)
	}
	again, err := w.service.Check(w.ctx, w.aigerim, attendance.CheckInput{Code: fresh.Code}, netip.Addr{})
	if err != nil || again.Status != attendance.Present {
		t.Fatalf("twice = %+v, %v", again, err)
	}
	if _, err := w.service.Check(w.ctx, w.aigerim, attendance.CheckInput{Code: "12a4"}, netip.Addr{}); !errors.Is(err, attendance.ErrInvalidCheck) {
		t.Fatalf("bad code shape: %v", err)
	}
	w.settings.set(attendance.QR, false)
	if _, err := w.service.Check(w.ctx, w.aigerim, attendance.CheckInput{SessionID: s.ID, Token: "x"}, netip.Addr{}); !errors.Is(err, attendance.ErrMethodOff) {
		t.Fatalf("qr off: %v", err)
	}
	open, _ = w.service.OpenLessons(w.ctx, w.aigerim)
	if len(open) != 1 || *open[0].Status != attendance.Present || open[0].Methods.QR {
		t.Fatalf("open after = %+v", open)
	}
	w.clock.move(2 * time.Hour)
	if open, _ := w.service.OpenLessons(w.ctx, w.aigerim); len(open) != 0 {
		t.Fatalf("ended lesson still open = %+v", open)
	}
}

func TestCampusNetworkAndGuessLimit(t *testing.T) {
	w := newWorld(t)
	l := w.lesson("2026-09-08", time.Date(2026, 9, 8, 9, 0, 0, 0, time.UTC))
	s, err := w.service.Open(w.ctx, w.teacher, l.SeriesID, l.Date)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	w.settings.set(attendance.Networks, "192.168.0.0/16, 172.16.0.1")
	if rec := w.call(w.aigerim, "POST", "/v1/attendance/check", `{"code":"`+s.Code+`"}`); rec.Code != 403 || !strings.Contains(rec.Body.String(), attendance.CodeOffCampus) {
		t.Fatalf("off campus = %d %s", rec.Code, rec.Body)
	}
	w.settings.set(attendance.Networks, "10.0.0.0/8")
	if rec := w.call(w.aigerim, "POST", "/v1/attendance/check", `{"code":"`+s.Code+`"}`); rec.Code != 200 || !strings.Contains(rec.Body.String(), `"present"`) {
		t.Fatalf("on campus = %d %s", rec.Code, rec.Body)
	}
	wrong := "0000"
	if s.Code == wrong {
		wrong = "1111"
	}
	for range 5 {
		if _, err := w.service.Check(w.ctx, w.erlan, attendance.CheckInput{Code: wrong}, netip.MustParseAddr("10.1.1.1")); !errors.Is(err, attendance.ErrWrongCode) {
			t.Fatalf("wrong: %v", err)
		}
	}
	if _, err := w.service.Check(w.ctx, w.erlan, attendance.CheckInput{Code: s.Code}, netip.MustParseAddr("10.1.1.1")); !errors.Is(err, attendance.ErrTooManyTries) {
		t.Fatalf("limit: %v", err)
	}
	w.clock.move(11 * time.Minute)
	if _, err := w.service.Check(w.ctx, w.erlan, attendance.CheckInput{Code: s.Code}, netip.MustParseAddr("10.1.1.1")); !errors.Is(err, attendance.ErrWrongCode) {
		t.Fatalf("after pause the code has expired: %v", err)
	}
}
