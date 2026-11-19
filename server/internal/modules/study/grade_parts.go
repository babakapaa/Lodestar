package study

type GradeParts struct {
	Defense *float64 `json:"defense,omitempty"`
	Map     *float64 `json:"map,omitempty"`
	Work    *float64 `json:"work,omitempty"`
}
