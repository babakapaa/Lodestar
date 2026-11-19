package module

type HelpSteps struct {
	Question  int `json:"question"`
	Direction int `json:"direction"`
	Source    int `json:"source"`
	Fragment  int `json:"fragment"`
	Example   int `json:"example"`
}
