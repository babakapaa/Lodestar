package module

import "time"

type ExternalWork struct {
	Title string
	Max   float64
	DueAt *time.Time
}
