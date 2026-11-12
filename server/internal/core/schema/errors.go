package schema

import "errors"

var ErrUnknownKey = errors.New("unknown key")

var ErrInvalidValue = errors.New("value does not match its schema")

var ErrBadDefinition = errors.New("definition is malformed or duplicated")
