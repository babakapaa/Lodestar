package lti

type SyncReport struct {
	Members  int `json:"members"`
	Enrolled int `json:"enrolled"`
	Works    int `json:"works"`
	Scores   int `json:"scores"`
}
