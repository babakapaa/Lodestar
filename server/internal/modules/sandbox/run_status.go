package sandbox

const (
	RunQueued   = "queued"
	RunStarting = "starting"
	RunRunning  = "running"
	RunStopping = "stopping"
	RunFinished = "finished"
	RunFailed   = "failed"
	RunStopped  = "stopped"
)

const (
	ReasonTimeLimit = "time_limit"
	ReasonNodeLost  = "node_lost"
	ReasonVanished  = "vanished"
)
