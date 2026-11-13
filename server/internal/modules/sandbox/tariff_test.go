package sandbox

import (
	"testing"
	"time"
)

func TestNightTariffWrapsMidnightInTenantZone(t *testing.T) {
	at := func(s string) time.Time {
		v, err := time.Parse(time.RFC3339, s)
		if err != nil {
			t.Fatal(err)
		}
		return v
	}
	utc := tariff{discount: 50, from: 23, to: 7, loc: time.UTC}
	if f := utc.factor(at("2026-11-13T22:00:00Z"), at("2026-11-14T02:00:00Z")); f != 0.625 {
		t.Fatalf("utc factor = %v", f)
	}
	if f := utc.factor(at("2026-11-13T10:00:00Z"), at("2026-11-13T12:30:00Z")); f != 1 {
		t.Fatalf("day factor = %v", f)
	}
	almaty, err := time.LoadLocation("Asia/Almaty")
	if err != nil {
		t.Skip("no zone data")
	}
	local := tariff{discount: 40, from: 23, to: 7, loc: almaty}
	if f := local.factor(at("2026-11-13T18:00:00Z"), at("2026-11-13T20:00:00Z")); f != 0.6 {
		t.Fatalf("almaty factor = %v", f)
	}
	if f := (tariff{discount: 50, from: 7, to: 7, loc: time.UTC}).factor(at("2026-11-13T00:00:00Z"), at("2026-11-13T05:00:00Z")); f != 1 {
		t.Fatalf("empty window = %v", f)
	}
}
