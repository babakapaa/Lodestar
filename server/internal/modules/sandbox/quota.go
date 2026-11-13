package sandbox

type Quota struct {
	Term     *Term       `json:"term"`
	Lines    []QuotaLine `json:"lines"`
	Hours    float64     `json:"hours"`
	Reserved float64     `json:"reserved"`
	Teams    []TeamPool  `json:"teams"`
}
