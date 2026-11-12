package search

type Hit struct {
	Kind     string  `json:"kind"`
	ID       string  `json:"id"`
	Title    string  `json:"title"`
	Subtitle string  `json:"subtitle"`
	Link     string  `json:"link"`
	Score    float64 `json:"score"`
}
