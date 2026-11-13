package league

type WorkWeight struct {
	Kind    string `json:"kind"`
	ID      string `json:"id"`
	Title   string `json:"title"`
	Project bool   `json:"project"`
	Weight  int    `json:"weight"`
}
