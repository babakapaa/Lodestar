package delivery

type job struct {
	id       int64
	userID   string
	channel  string
	item     Item
	attempts int
}
