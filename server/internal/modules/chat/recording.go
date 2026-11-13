package chat

type Recording struct {
	EgressID  string
	CallID    string
	ChatID    string
	FileID    *string
	StartedBy *string
	AudioOnly bool
}
