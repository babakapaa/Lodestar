package module

type LinkTarget struct {
	Kind  string `json:"kind"`
	ID    string `json:"id"`
	Title string `json:"title"`
	Path  string `json:"path"`
}
