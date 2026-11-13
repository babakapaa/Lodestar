package league

import "time"

type League struct {
	ID         string
	OriginKind string
	OriginID   string
	Divisions  Divisions
	Confirm    bool
	UpdatedAt  time.Time
}
