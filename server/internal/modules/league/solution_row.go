package league

type SolutionRow struct {
	Solution
	Similar *SimilarSolution `json:"similar"`
}
