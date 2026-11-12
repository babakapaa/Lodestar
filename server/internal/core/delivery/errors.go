package delivery

import "errors"

var ErrInvalidDevice = errors.New("device platform, token, name or language is invalid")

var ErrDeviceNotFound = errors.New("device not found")
