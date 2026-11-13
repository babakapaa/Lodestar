package league

import "time"

type Dispute struct {
	ID           string        `json:"id"`
	UserID       string        `json:"userId"`
	EntryID      *string       `json:"entryId"`
	TransitionID *string       `json:"transitionId"`
	Text         string        `json:"text"`
	Status       DisputeStatus `json:"status"`
	Answer       string        `json:"answer"`
	DecidedAt    *time.Time    `json:"decidedAt"`
	CreatedAt    time.Time     `json:"createdAt"`
}
