package audit

type Filter struct {
	Before  int64
	Action  string
	ActorID string
}
