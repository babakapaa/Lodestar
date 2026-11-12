package notify

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type PrivacySection struct {
	store    NotificationStore
	channels ChannelStore
}

func (PrivacySection) Name() string {
	return "notifications"
}

func (s PrivacySection) Export(ctx context.Context, tx pgx.Tx, userID string) (any, error) {
	list, err := s.store.All(ctx, tx, userID)
	if err != nil {
		return nil, err
	}
	if list == nil {
		list = []Notification{}
	}
	channels, err := s.channels.ForUser(ctx, tx, userID)
	if err != nil {
		return nil, err
	}
	return struct {
		Notifications []Notification           `json:"notifications"`
		Channels      map[string]ChannelsInput `json:"channels"`
	}{list, channels}, nil
}
