package navigator

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type PrivacySection struct {
	service Service
}

func (PrivacySection) Name() string {
	return "navigator.dialogs"
}

func (s PrivacySection) Export(ctx context.Context, tx pgx.Tx, userID string) (any, error) {
	dialogs, err := s.service.dialogs.List(ctx, tx, userID, "", "")
	if err != nil {
		return nil, err
	}
	out := make([]map[string]any, 0, len(dialogs))
	for _, d := range dialogs {
		messages, err := s.service.messages.List(ctx, tx, d.ID)
		if err != nil {
			return nil, err
		}
		out = append(out, map[string]any{"dialog": d, "messages": messages})
	}
	return out, nil
}
