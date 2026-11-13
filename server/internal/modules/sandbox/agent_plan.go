package sandbox

type AgentPlan struct {
	NodeID string    `json:"nodeId"`
	Runs   []RunSpec `json:"runs"`
}
