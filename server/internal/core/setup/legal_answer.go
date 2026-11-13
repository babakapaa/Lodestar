package setup

import "time"

type LegalAnswer struct {
	Question   string     `json:"question"`
	Title      string     `json:"title"`
	Answer     string     `json:"answer"`
	Basis      string     `json:"basis"`
	AnsweredBy *string    `json:"answeredBy"`
	AnsweredAt *time.Time `json:"answeredAt"`
	ApprovedBy *string    `json:"approvedBy"`
	ApprovedAt *time.Time `json:"approvedAt"`
}
