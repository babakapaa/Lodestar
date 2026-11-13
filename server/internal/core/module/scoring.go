package module

type Scoring interface {
	ScoreSources(k Kernel) []ScoreSource
}
