package league

type RulesInput struct {
	Divisions Divisions     `json:"divisions"`
	Confirm   bool          `json:"confirm"`
	Weights   []WeightInput `json:"weights"`
}

func (in RulesInput) normalized(works []WorkWeight, gpuCap int) (RulesInput, error) {
	divisions, err := in.Divisions.normalized(gpuCap)
	if err != nil {
		return in, err
	}
	in.Divisions = divisions
	known := make(map[string]bool, len(works))
	for _, w := range works {
		known[w.Kind+"/"+w.ID] = true
	}
	seen := map[string]bool{}
	for _, w := range in.Weights {
		key := w.Kind + "/" + w.ID
		if !known[key] || seen[key] || w.Weight < 0 || w.Weight > maxWeight {
			return in, ErrInvalidRules
		}
		seen[key] = true
	}
	return in, nil
}
