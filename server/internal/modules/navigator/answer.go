package navigator

type Answer struct {
	Question Message `json:"question"`
	Reply    Message `json:"reply"`
}
