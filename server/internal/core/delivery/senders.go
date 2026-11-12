package delivery

import (
	"context"

	"github.com/babakapaa/lodestar/server/internal/platform/fcm"
	"github.com/babakapaa/lodestar/server/internal/platform/mail"
)

type Pusher interface {
	Push(ctx context.Context, token string, m fcm.Message) error
}

type Mailer interface {
	Send(ctx context.Context, m mail.Message) error
}
