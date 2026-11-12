package fcm

import "errors"

var ErrUnregistered = errors.New("device token is no longer registered")
