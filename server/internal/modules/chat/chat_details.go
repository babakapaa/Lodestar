package chat

type ChatDetails struct {
	Chat
	Members     []Member `json:"members"`
	PeerReadSeq *int64   `json:"peerReadSeq"`
}
