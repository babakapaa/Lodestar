package sandbox

const (
	SourceBase     = "sandbox.base"
	SourceTopup    = "sandbox.topup"
	SourceTransfer = "sandbox.transfer"
	SourceRun      = "sandbox.run"
)

type QuotaLine struct {
	Source string  `json:"source"`
	Title  string  `json:"title"`
	Hours  float64 `json:"hours"`
}
