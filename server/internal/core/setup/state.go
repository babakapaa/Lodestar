package setup

import "time"

type State struct {
	ProtocolVersion int
	SignedBy        *string
	SignedAt        *time.Time
	LaunchedAt      *time.Time
	Pilot           []string
}
