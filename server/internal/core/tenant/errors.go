package tenant

import "errors"

var ErrNotFound = errors.New("tenant not found")

var ErrSlugTaken = errors.New("tenant slug already taken")

var ErrInvalid = errors.New("tenant slug or name is invalid")
