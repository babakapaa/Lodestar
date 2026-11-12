package study

type MapRules struct {
	MinSources   int `json:"minSources"`
	MinScholarly int `json:"minScholarly"`
}

func (r MapRules) valid() bool {
	return r.MinSources >= 0 && r.MinSources <= maxMapSources && r.MinScholarly >= 0 && r.MinScholarly <= 100
}
