package sandbox

const (
	PhaseRunning = "running"
	PhaseExited  = "exited"
)

type RunReport struct {
	RunID    string    `json:"runId"`
	Phase    string    `json:"phase"`
	URL      string    `json:"url"`
	ExitCode int       `json:"exitCode"`
	Usage    *RunUsage `json:"usage"`
	Log      string    `json:"log"`
}
