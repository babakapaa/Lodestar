package study

import "math"

type GradeWeights struct {
	Defense int `json:"defense"`
	Map     int `json:"map"`
	Work    int `json:"work"`
}

func (w GradeWeights) valid() bool {
	return w.Defense >= 0 && w.Map >= 0 && w.Work >= 0 && w.Defense+w.Map+w.Work == 100
}

func (w GradeWeights) points(parts GradeParts, maxPoints float64) (float64, bool) {
	share := 0.0
	for _, p := range []struct {
		weight int
		value  *float64
	}{{w.Defense, parts.Defense}, {w.Map, parts.Map}, {w.Work, parts.Work}} {
		if p.weight == 0 {
			if p.value != nil {
				return 0, false
			}
			continue
		}
		if p.value == nil || *p.value < 0 || *p.value > 100 {
			return 0, false
		}
		share += float64(p.weight) * *p.value
	}
	return math.Round(share/10000*maxPoints*100) / 100, true
}
