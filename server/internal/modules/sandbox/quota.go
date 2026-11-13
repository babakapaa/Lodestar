package sandbox

type Quota struct {
	Term     *Term       `json:"term"`
	Lines    []QuotaLine `json:"lines"`
	Hours    float64     `json:"hours"`
	Reserved float64     `json:"reserved"`
	Spent    float64     `json:"spent"`
	Forecast float64     `json:"forecast"`
	Teams    []TeamPool  `json:"teams"`
}
