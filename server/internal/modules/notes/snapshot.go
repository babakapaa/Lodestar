package notes

type Snapshot struct {
	Seq  int64  `json:"seq"`
	Data []byte `json:"data"`
}
