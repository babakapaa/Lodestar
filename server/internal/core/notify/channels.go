package notify

type Channels struct {
	Type     string `json:"type"`
	Required bool   `json:"required"`
	InApp    bool   `json:"inApp"`
	Push     bool   `json:"push"`
	Email    bool   `json:"email"`
}

type ChannelsInput struct {
	InApp bool `json:"inApp"`
	Push  bool `json:"push"`
	Email bool `json:"email"`
}
