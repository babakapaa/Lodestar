package study

import (
	"errors"

	"github.com/babakapaa/lodestar/server/internal/core/access"
)

func isForbidden(err error) bool {
	return errors.Is(err, access.ErrForbidden)
}
