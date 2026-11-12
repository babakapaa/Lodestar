package study

type Grading string

const (
	GradingPoints Grading = "points"
	GradingPass   Grading = "pass"
	GradingNone   Grading = "none"
)

func (g Grading) valid() bool {
	return g == GradingPoints || g == GradingPass || g == GradingNone
}
