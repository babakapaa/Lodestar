package livekit

const (
	EgressEnded    = "egress_ended"
	EgressComplete = "EGRESS_COMPLETE"
)

type EgressInfo struct {
	EgressID    string       `json:"egressId"`
	RoomName    string       `json:"roomName"`
	Status      string       `json:"status"`
	FileResults []FileResult `json:"fileResults"`
}
