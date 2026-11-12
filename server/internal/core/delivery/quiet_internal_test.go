package delivery

import (
	"testing"
	"time"

	"github.com/babakapaa/lodestar/server/internal/core/module"
)

func TestQuietHoursOvernight(t *testing.T) {
	q := quietHours{on: true, from: 22 * 60, to: 7 * 60, location: time.UTC}
	cases := []struct {
		at    time.Time
		quiet bool
		until time.Time
	}{
		{time.Date(2026, 1, 1, 21, 59, 0, 0, time.UTC), false, time.Time{}},
		{time.Date(2026, 1, 1, 22, 0, 0, 0, time.UTC), true, time.Date(2026, 1, 2, 7, 0, 0, 0, time.UTC)},
		{time.Date(2026, 1, 2, 3, 0, 0, 0, time.UTC), true, time.Date(2026, 1, 2, 7, 0, 0, 0, time.UTC)},
		{time.Date(2026, 1, 2, 7, 0, 0, 0, time.UTC), false, time.Time{}},
	}
	for _, c := range cases {
		until, quiet := q.until(c.at)
		if quiet != c.quiet || !until.Equal(c.until) {
			t.Errorf("until(%v) = %v, %v", c.at, until, quiet)
		}
	}
	if _, quiet := (quietHours{from: 0, to: 600, location: time.UTC}).until(time.Date(2026, 1, 1, 3, 0, 0, 0, time.UTC)); quiet {
		t.Error("switched off quiet hours still apply")
	}
}

func TestTextsFallBack(t *testing.T) {
	texts := NewTexts(module.NotificationType{Key: "a.b", Texts: map[string]module.NotificationText{
		"ru": {Title: "ру"}, "kk": {Title: "қаз"},
	}})
	item := Item{Type: "a.b", Params: []byte(`{}`)}
	if got, _ := texts.Render(item, "", "fr", "kk"); got.Title != "қаз" {
		t.Errorf("chosen = %q", got.Title)
	}
	if got, _ := texts.Render(item, "de"); got.Title != "қаз" {
		t.Errorf("fallback = %q, want first by code", got.Title)
	}
	if _, ok := texts.Render(Item{Type: "x.y"}); ok {
		t.Error("unknown type rendered")
	}
}
