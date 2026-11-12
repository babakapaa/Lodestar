package app

import (
	"time"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/platform/config"
)

func Lifetimes(cfg config.Config) auth.Lifetimes {
	return auth.Lifetimes{
		Access:  cfg.AccessTTL,
		Refresh: cfg.RefreshTTL,
		Login:   10 * time.Minute,
		Grant:   time.Minute,
	}
}
