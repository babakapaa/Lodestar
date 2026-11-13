package sandbox

type NodeCreated struct {
	Node
	Token string `json:"token"`
}
