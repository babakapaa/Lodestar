package study_test

import (
	"testing"
	"time"

	"github.com/babakapaa/lodestar/server/internal/modules/study"
)

func day(t *testing.T, value string) study.Day {
	t.Helper()
	d, err := study.ParseDay(value)
	if err != nil {
		t.Fatalf("day %s: %v", value, err)
	}
	return d
}

func TestCalendarExpandsEveryOtherMonday(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Almaty")
	s := study.Series{
		ID: "s1", CourseID: "c1", Kind: "lecture", Room: "301",
		StartsOn: day(t, "2026-09-02"), EndsOn: day(t, "2026-10-31"),
		Weekday: 1, StartTime: 9*60 + 30, Minutes: 80, EveryWeeks: 2,
	}
	moved := time.Date(2026, 9, 22, 14, 0, 0, 0, loc)
	room := "105"
	changes := []study.Change{
		{SeriesID: "s1", Date: day(t, "2026-09-07"), Status: study.ChangeCancelled, Note: "Праздник"},
		{SeriesID: "s1", Date: day(t, "2026-09-21"), Status: study.ChangeChanged, StartsAt: &moved, Room: &room},
	}
	lessons := study.Calendar{Location: loc}.Expand([]study.Series{s}, changes, map[string]string{"c1": "Сети"}, day(t, "2026-09-01"), day(t, "2026-10-10"))

	var dates []string
	for _, l := range lessons {
		dates = append(dates, l.Date.String())
	}
	want := []string{"2026-09-07", "2026-09-21", "2026-10-05"}
	if len(dates) != len(want) {
		t.Fatalf("dates = %v, want %v", dates, want)
	}
	for i := range want {
		if dates[i] != want[i] {
			t.Fatalf("dates = %v, want %v", dates, want)
		}
	}
	if lessons[0].Status != study.LessonCancelled || lessons[0].Note != "Праздник" {
		t.Fatalf("cancelled = %+v", lessons[0])
	}
	if lessons[1].Status != study.LessonChanged || !lessons[1].StartsAt.Equal(moved) || lessons[1].Room != "105" || lessons[1].EndsAt.Sub(lessons[1].StartsAt) != 80*time.Minute {
		t.Fatalf("moved = %+v", lessons[1])
	}
	if got := lessons[2].StartsAt; got.Hour() != 9 || got.Minute() != 30 || got.Location() != loc || lessons[2].CourseTitle != "Сети" {
		t.Fatalf("planned = %+v", lessons[2])
	}
	if !(study.Calendar{}).Occurs(s, day(t, "2026-10-19")) || (study.Calendar{}).Occurs(s, day(t, "2026-10-12")) {
		t.Fatal("occurs is wrong")
	}
}
