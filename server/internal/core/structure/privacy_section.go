package structure

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type PrivacySection struct {
	members MemberStore
}

func (PrivacySection) Name() string {
	return "units"
}

func (s PrivacySection) Export(ctx context.Context, tx pgx.Tx, userID string) (any, error) {
	list, err := s.members.UnitsOf(ctx, tx, userID)
	if list == nil {
		list = []Unit{}
	}
	return list, err
}
