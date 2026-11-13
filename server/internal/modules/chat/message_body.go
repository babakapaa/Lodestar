package chat

type messageBody struct {
	Text     string   `json:"text"`
	Mentions []string `json:"mentions,omitempty"`
}
