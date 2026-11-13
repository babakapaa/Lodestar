package league

import (
	"slices"
	"time"
)

type Solution struct {
	ID          string         `json:"id"`
	TaskID      string         `json:"taskId"`
	TeamID      string         `json:"teamId"`
	TeamName    string         `json:"teamName"`
	ProposedBy  *string        `json:"proposedBy"`
	Members     []string       `json:"members"`
	Work        SolutionWork   `json:"work"`
	Votes       []SolutionVote `json:"votes"`
	SimilarTo   *string        `json:"similarTo"`
	Similarity  *float64       `json:"similarity"`
	Explanation string         `json:"explanation"`
	Points      *int           `json:"points"`
	Passed      []string       `json:"passed"`
	Comment     string         `json:"comment"`
	ReviewedAt  *time.Time     `json:"reviewedAt"`
	CreatedAt   time.Time      `json:"createdAt"`
}

func (s Solution) has(userID string) bool {
	return slices.Contains(s.Members, userID)
}
