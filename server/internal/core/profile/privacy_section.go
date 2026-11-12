package profile

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
)

type PrivacySection struct {
	profiles    ProfileStore
	preferences PreferenceStore
}

func (PrivacySection) Name() string {
	return "profile"
}

func (s PrivacySection) Export(ctx context.Context, tx pgx.Tx, userID string) (any, error) {
	p, err := s.profiles.Get(ctx, tx, userID)
	if err != nil {
		return nil, err
	}
	prefs, err := s.preferences.All(ctx, tx, userID)
	if err != nil {
		return nil, err
	}
	return struct {
		Profile     Profile                    `json:"profile"`
		Preferences map[string]json.RawMessage `json:"preferences"`
	}{p, prefs}, nil
}
