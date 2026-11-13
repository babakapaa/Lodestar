package module

type ScoredWork struct {
	Kind    string
	ID      string
	Title   string
	Project bool
	Max     float64
	Results []WorkResult
}
