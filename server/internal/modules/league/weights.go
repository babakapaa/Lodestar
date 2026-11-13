package league

import "github.com/babakapaa/lodestar/server/internal/core/module"

const (
	maxWeight     = 10000
	projectWeight = 100
)

type Weights map[string]int

func (w Weights) of(work module.ScoredWork) int {
	if weight, ok := w[work.Kind+"/"+work.ID]; ok {
		return weight
	}
	if work.Project {
		return projectWeight
	}
	return 0
}

func (w Weights) formula(works []module.ScoredWork) []WorkWeight {
	out := make([]WorkWeight, 0, len(works))
	for _, work := range works {
		out = append(out, WorkWeight{Kind: work.Kind, ID: work.ID, Title: work.Title, Project: work.Project, Weight: w.of(work)})
	}
	return out
}
