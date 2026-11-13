package chat

type MessagePage struct {
	Messages []Message `json:"messages"`
	More     bool      `json:"more"`
}
