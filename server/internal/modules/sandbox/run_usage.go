package sandbox

type RunUsage struct {
	CPU       float64  `json:"cpu"`
	GPU       *float64 `json:"gpu"`
	NetBytes  int64    `json:"netBytes"`
	Processes []string `json:"processes"`
}

const (
	activeCPU   = 5
	activeGPU   = 5
	activeBytes = 256 << 10
)

func (u RunUsage) active(lastNet int64) bool {
	return u.CPU >= activeCPU || (u.GPU != nil && *u.GPU >= activeGPU) || u.NetBytes-lastNet >= activeBytes
}
