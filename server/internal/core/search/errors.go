package search

import "errors"

var ErrInvalidQuery = errors.New("search query is empty or too long")

var ErrUnknownKind = errors.New("unknown search kind")

var ErrBadSource = errors.New("search source kind is malformed or duplicated")
