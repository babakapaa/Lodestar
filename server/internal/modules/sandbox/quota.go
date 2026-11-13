package sandbox

type Quota struct {
	Term  *Term       `json:"term"`
	Lines []QuotaLine `json:"lines"`
	Hours float64     `json:"hours"`
	Teams []TeamPool  `json:"teams"`
}
