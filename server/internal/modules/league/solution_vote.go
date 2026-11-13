package league

type SolutionVote struct {
	UserID   string `json:"userId"`
	Decision string `json:"decision"`
	Comment  string `json:"comment"`
}
