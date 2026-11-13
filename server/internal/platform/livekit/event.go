package livekit

const (
	ParticipantLeft = "participant_left"
	RoomFinished    = "room_finished"
)

type Event struct {
	Event       string      `json:"event"`
	Room        Room        `json:"room"`
	Participant Participant `json:"participant"`
}
