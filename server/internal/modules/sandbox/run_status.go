package sandbox

const (
	RunQueued   = "queued"
	RunStarting = "starting"
	RunRunning  = "running"
	RunHeld     = "held"
	RunStopping = "stopping"
	RunFinished = "finished"
	RunFailed   = "failed"
	RunStopped  = "stopped"
)

const (
	ReasonTimeLimit = "time_limit"
	ReasonNodeLost  = "node_lost"
	ReasonVanished  = "vanished"
	ReasonIdle      = "idle"
	ReasonWatchdog  = "watchdog"
)
