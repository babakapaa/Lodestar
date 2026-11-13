package navigator

type checkReply struct {
	Unsupported []int `json:"unsupported"`
	Leak        bool  `json:"leak"`
}
