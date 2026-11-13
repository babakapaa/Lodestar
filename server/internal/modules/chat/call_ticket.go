package chat

type CallTicket struct {
	Call        Call   `json:"call"`
	URL         string `json:"url"`
	Room        string `json:"room"`
	Token       string `json:"token"`
	ScreenShare bool   `json:"screenShare"`
	Record      bool   `json:"record"`
	Muted       bool   `json:"muted"`
}
