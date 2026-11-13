package setup

import "time"

type Protocol struct {
	Version  int           `json:"version"`
	SignedBy *string       `json:"signedBy"`
	SignedAt *time.Time    `json:"signedAt"`
	Answers  []LegalAnswer `json:"answers"`
}
