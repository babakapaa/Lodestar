package livekit

type videoClaims struct {
	Room              string   `json:"room,omitempty"`
	RoomJoin          bool     `json:"roomJoin,omitempty"`
	RoomRecord        bool     `json:"roomRecord,omitempty"`
	CanPublish        bool     `json:"canPublish,omitempty"`
	CanSubscribe      bool     `json:"canSubscribe,omitempty"`
	CanPublishData    bool     `json:"canPublishData,omitempty"`
	CanPublishSources []string `json:"canPublishSources,omitempty"`
}
