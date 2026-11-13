package sandbox

const (
	PhaseRunning = "running"
	PhaseExited  = "exited"
)

type RunReport struct {
	RunID    string `json:"runId"`
	Phase    string `json:"phase"`
	URL      string `json:"url"`
	ExitCode int    `json:"exitCode"`
}
