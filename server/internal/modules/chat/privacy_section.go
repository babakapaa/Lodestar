package chat

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type PrivacySection struct {
	messages MessageStore
}

func (PrivacySection) Name() string {
	return "chat.messages"
}

func (s PrivacySection) Export(ctx context.Context, tx pgx.Tx, userID string) (any, error) {
	return s.messages.ByAuthor(ctx, tx, userID)
}
