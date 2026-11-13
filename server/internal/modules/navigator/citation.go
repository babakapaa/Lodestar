package navigator

type Citation struct {
	N          int    `json:"n"`
	MaterialID string `json:"materialId"`
	Title      string `json:"title"`
	Reference  string `json:"reference"`
	URL        string `json:"url"`
	Text       string `json:"text"`
}
