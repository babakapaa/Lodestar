package notes

type Changes struct {
	Seq      int64     `json:"seq"`
	Snapshot *Snapshot `json:"snapshot"`
	Updates  []Update  `json:"updates"`
}
