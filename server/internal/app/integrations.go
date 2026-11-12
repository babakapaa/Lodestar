package app

import (
	"github.com/babakapaa/lodestar/server/internal/core/delivery"
	"github.com/babakapaa/lodestar/server/internal/core/files"
)

type Integrations struct {
	Storage files.Storage
	Pusher  delivery.Pusher
	Mailer  delivery.Mailer
}
