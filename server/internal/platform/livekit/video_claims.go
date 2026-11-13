package livekit

type videoClaims struct {
	Room              string   `json:"room"`
	RoomJoin          bool     `json:"roomJoin"`
	CanPublish        bool     `json:"canPublish"`
	CanSubscribe      bool     `json:"canSubscribe"`
	CanPublishData    bool     `json:"canPublishData"`
	CanPublishSources []string `json:"canPublishSources"`
}
