package app

import (
	"github.com/babakapaa/lodestar/server/internal/core/delivery"
	"github.com/babakapaa/lodestar/server/internal/core/files"
	"github.com/babakapaa/lodestar/server/internal/platform/llm"
)

type Integrations struct {
	Storage files.Storage
	Pusher  delivery.Pusher
	Mailer  delivery.Mailer
	Model   llm.Model
}
