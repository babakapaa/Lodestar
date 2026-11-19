package navigator

type messageBody struct {
	Text      string     `json:"text"`
	Intent    Intent     `json:"intent,omitempty"`
	Step      Step       `json:"step,omitempty"`
	Support   string     `json:"support,omitempty"`
	Citations []Citation `json:"citations,omitempty"`
	Articles  []Article  `json:"articles,omitempty"`
}
