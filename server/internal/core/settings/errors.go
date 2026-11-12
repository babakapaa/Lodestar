package settings

import "errors"

var ErrUnknownModule = errors.New("module not found")

var ErrRequiredModule = errors.New("module cannot be turned off")
