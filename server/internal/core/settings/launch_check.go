package settings

import (
	"context"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
)

type LaunchCheck interface {
	Allowed(ctx context.Context, p auth.Principal) (bool, error)
}
