package live

type EventsPage struct {
	Events []Event `json:"events"`
	More   bool    `json:"more"`
}
