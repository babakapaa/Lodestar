package privacy

import "time"

type manifest struct {
	UserID     string    `json:"userId"`
	ExportedAt time.Time `json:"exportedAt"`
	Sections   []string  `json:"sections"`
	Files      []string  `json:"files"`
}
