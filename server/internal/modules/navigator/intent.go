package navigator

import "slices"

type Intent string

const (
	IntentUnderstand Intent = "understand"
	IntentFind       Intent = "find"
	IntentCheck      Intent = "check"
	IntentDoForMe    Intent = "do_for_me"
	IntentOrganize   Intent = "organize"
	IntentState      Intent = "state"
	IntentClosed     Intent = "closed"
)

var chosenIntents = []Intent{IntentUnderstand, IntentFind, IntentCheck, IntentDoForMe, IntentOrganize, IntentState}

func (i Intent) Choosable() bool {
	return slices.Contains(chosenIntents, i)
}

func (i Intent) start() Step {
	switch i {
	case IntentFind:
		return StepSource
	case IntentUnderstand, IntentCheck, IntentDoForMe:
		return StepQuestion
	default:
		return StepNone
	}
}
