package module

type Notification struct {
	Type    string
	UserIDs []string
	Params  map[string]any
	Link    string
}
