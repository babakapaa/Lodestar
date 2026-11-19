package navigator

import "slices"

type Step string

const (
	StepNone      Step = ""
	StepQuestion  Step = "question"
	StepDirection Step = "direction"
	StepSource    Step = "source"
	StepFragment  Step = "fragment"
	StepExample   Step = "example"
)

var ladder = []Step{StepQuestion, StepDirection, StepSource, StepFragment, StepExample}

func (s Step) next() Step {
	i := slices.Index(ladder, s)
	if i < 0 {
		return StepQuestion
	}
	return ladder[min(i+1, len(ladder)-1)]
}

func (s Step) cites() bool {
	return s == StepSource || s == StepFragment || s == StepExample
}

func (s Step) Last() bool {
	return s == StepExample
}
