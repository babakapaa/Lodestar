package navigator

type TopicCount struct {
	MaterialID string `json:"materialId"`
	Title      string `json:"title"`
	Students   int    `json:"students"`
	Recent     int    `json:"recent"`
}
