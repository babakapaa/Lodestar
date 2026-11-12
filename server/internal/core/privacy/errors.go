package privacy

import "errors"

var ErrBadSection = errors.New("privacy section name is malformed or duplicated")

var ErrSelfDeletionOff = errors.New("the university handles account deletion itself")

var ErrNoRequest = errors.New("no deletion request")
