package setup

import "errors"

var (
	ErrUnknownSection  = errors.New("unknown setup section")
	ErrUnknownQuestion = errors.New("unknown legal question")
	ErrWrongSection    = errors.New("setting belongs to another section")
	ErrInvalidAnswer   = errors.New("answer and basis are required")
	ErrInvalidAssignee = errors.New("assignee cannot take this section")
	ErrInvalidStatus   = errors.New("section is not in a state for this step")
	ErrLegalOpen       = errors.New("legal answers are not all approved")
	ErrNotReady        = errors.New("setup has open problems")
	ErrNoSetup         = errors.New("university is not in setup")
)
