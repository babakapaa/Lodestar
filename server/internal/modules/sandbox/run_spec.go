package sandbox

type RunSpec struct {
	RunID     string  `json:"runId"`
	Kind      string  `json:"kind"`
	Title     string  `json:"title"`
	Image     string  `json:"image"`
	Command   string  `json:"command"`
	MaxHours  float64 `json:"maxHours"`
	Workspace string  `json:"workspace"`
	Hold      bool    `json:"hold"`
	Stop      bool    `json:"stop"`
}
