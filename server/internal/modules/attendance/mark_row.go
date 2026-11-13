package attendance

import "time"

type markRow struct {
	SessionID string
	UserID    string
	Status    *MarkStatus
	Method    *MarkMethod
	MarkedAt  *time.Time
}
