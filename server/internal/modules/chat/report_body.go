package chat

type reportBody struct {
	ChatTitle string            `json:"chatTitle"`
	Reason    string            `json:"reason"`
	Messages  []ReportedMessage `json:"messages"`
}
