package agent

import "strconv"

const (
	labelRun   = "lodestar.run"
	labelGPU   = "lodestar.gpu"
	labelToken = "lodestar.token"
	labelKind  = "lodestar.kind"
)

type Container struct {
	ID       string
	RunID    string
	Kind     string
	GPU      int
	Token    string
	State    string
	ExitCode int
	Port     int
}

func (c Container) running() bool {
	return c.State == "running"
}

func (c Container) paused() bool {
	return c.State == "paused"
}

func (c Container) exited() bool {
	return c.State == "exited" || c.State == "dead"
}

func fromLabels(id, state string, labels map[string]string) Container {
	gpu, err := strconv.Atoi(labels[labelGPU])
	if err != nil {
		gpu = -1
	}
	return Container{ID: id, RunID: labels[labelRun], Kind: labels[labelKind], GPU: gpu, Token: labels[labelToken], State: state}
}
