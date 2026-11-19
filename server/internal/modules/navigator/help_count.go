package navigator

import (
	"slices"

	"github.com/babakapaa/lodestar/server/internal/core/module"
)

const maxHelpSources = 50

func emptyHelp() module.HelpSummary {
	return module.HelpSummary{Sources: []string{}, Articles: []string{}}
}

func filledHelp(sum module.HelpSummary) module.HelpSummary {
	if sum.Sources == nil {
		sum.Sources = []string{}
	}
	if sum.Articles == nil {
		sum.Articles = []string{}
	}
	return sum
}

func countHelp(sum module.HelpSummary, body messageBody) module.HelpSummary {
	sum = filledHelp(sum)
	sum.Replies++
	switch body.Step {
	case StepQuestion:
		sum.Steps.Question++
	case StepDirection:
		sum.Steps.Direction++
	case StepSource:
		sum.Steps.Source++
	case StepFragment:
		sum.Steps.Fragment++
	case StepExample:
		sum.Steps.Example++
	}
	for _, c := range body.Citations {
		sum.Sources = addOnce(sum.Sources, c.label())
	}
	for _, a := range body.Articles {
		sum.Articles = addOnce(sum.Articles, a.label())
	}
	return sum
}

func addOnce(list []string, item string) []string {
	if item == "" || len(list) >= maxHelpSources || slices.Contains(list, item) {
		return list
	}
	return append(list, item)
}
