package sandbox

import "math"

func hours(v float64) float64 {
	return math.Round(v*100) / 100
}
