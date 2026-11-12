package notify

type Message struct {
	Type    string
	UserIDs []string
	Params  map[string]any
	Link    string
}
