package navigator

type Classification struct {
	Intent   Intent `json:"intent"`
	Distress bool   `json:"distress"`
}
