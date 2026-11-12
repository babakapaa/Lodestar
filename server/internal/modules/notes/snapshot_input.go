package notes

const maxSnapshotBytes = 64 << 20

type SnapshotInput struct {
	Seq  int64  `json:"seq"`
	Data []byte `json:"data"`
	Body string `json:"body"`
}

func (in SnapshotInput) valid() bool {
	return in.Seq > 0 && len(in.Data) > 0 && len(in.Data)+len(in.Body) <= maxSnapshotBytes
}
