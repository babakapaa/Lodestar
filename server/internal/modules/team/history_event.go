package team

type HistoryEvent string

const (
	EventCreated   HistoryEvent = "created"
	EventJoined    HistoryEvent = "joined"
	EventAdded     HistoryEvent = "added"
	EventLeft      HistoryEvent = "left"
	EventRemoved   HistoryEvent = "removed"
	EventSeparated HistoryEvent = "separated"
	EventMovedIn   HistoryEvent = "moved_in"
	EventMovedOut  HistoryEvent = "moved_out"
)
