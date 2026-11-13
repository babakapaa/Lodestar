package sandbox

import "time"

type runEnd struct {
	status   string
	exitCode *int
	reason   string
	charged  float64
	refunded bool
	log      string
	at       time.Time
}
