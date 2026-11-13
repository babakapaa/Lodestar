package sandbox

type TeamPool struct {
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	Hours    float64     `json:"hours"`
	Reserved float64     `json:"reserved"`
	Lines    []QuotaLine `json:"lines"`
}
