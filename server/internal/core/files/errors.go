package files

import "errors"

var ErrInvalidFile = errors.New("file name, type or size is invalid")

var ErrTooLarge = errors.New("file is larger than the university allows")

var ErrNotFound = errors.New("file not found")

var ErrNotUploaded = errors.New("file content has not been uploaded")

var ErrSizeMismatch = errors.New("uploaded content does not match the declared size")

var ErrUnavailable = errors.New("file storage is not configured")
