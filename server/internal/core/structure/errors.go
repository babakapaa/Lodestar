package structure

import "errors"

var ErrUnitNotFound = errors.New("unit not found")

var ErrInvalidUnit = errors.New("unit kind, name or code is invalid")

var ErrCodeTaken = errors.New("unit code already taken")

var ErrCycle = errors.New("unit cannot be moved inside itself")

var ErrArchived = errors.New("unit is archived")

var ErrUserNotFound = errors.New("user not found")

var ErrBadImport = errors.New("import file is malformed")
