package chat

type ConsentRequest struct {
	RequestedBy string   `json:"requestedBy"`
	AudioOnly   bool     `json:"audioOnly"`
	Waiting     []string `json:"waiting"`
}
