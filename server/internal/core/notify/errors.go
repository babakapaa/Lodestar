package notify

import "errors"

var ErrUnknownType = errors.New("unknown notification type")

var ErrBadType = errors.New("notification type is malformed or duplicated")

var ErrRequiredType = errors.New("required notification cannot be turned off in the app")

var ErrNotFound = errors.New("notification not found")
