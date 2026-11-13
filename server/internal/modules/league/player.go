package league

import "time"

type Player struct {
	UserID   string
	Division int
	Display  Display
	Nickname string
	JoinedAt time.Time
}
