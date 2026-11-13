package navigator

import "errors"

var ErrUnavailable = errors.New("navigator model is not configured")

var ErrDialogNotFound = errors.New("dialog not found")

var ErrInvalidDialog = errors.New("dialog id or scope is invalid")

var ErrInvalidMessage = errors.New("message is empty, too long or asks for a closed type")

var ErrModelFailed = errors.New("navigator could not get an answer from the model")

var ErrInvalidMaterial = errors.New("material title, link or text is empty or too long")

var ErrMaterialNotFound = errors.New("material not found")

var ErrJournalNotFound = errors.New("journal entry not found")
