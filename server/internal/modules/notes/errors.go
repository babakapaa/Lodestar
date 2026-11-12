package notes

import "errors"

var (
	ErrNoteNotFound    = errors.New("note not found")
	ErrInvalidNote     = errors.New("note is invalid")
	ErrNoteExists      = errors.New("note id belongs to someone else")
	ErrReadOnly        = errors.New("note is read only")
	ErrNoteDeleted     = errors.New("note is deleted")
	ErrInvalidUpdate   = errors.New("note update is invalid")
	ErrStaleSnapshot   = errors.New("snapshot does not match the note")
	ErrQuotaExceeded   = errors.New("notes quota exceeded")
	ErrUserNotFound    = errors.New("user not found")
	ErrNotOwner        = errors.New("only the owner can do this")
	ErrVersionNotFound = errors.New("note version not found")
)
