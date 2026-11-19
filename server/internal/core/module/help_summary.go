package module

type HelpSummary struct {
	Replies  int       `json:"replies"`
	Steps    HelpSteps `json:"steps"`
	Sources  []string  `json:"sources"`
	Articles []string  `json:"articles"`
}
