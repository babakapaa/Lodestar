package delivery

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type PrivacySection struct {
	devices DeviceStore
}

func (PrivacySection) Name() string {
	return "devices"
}

func (s PrivacySection) Export(ctx context.Context, tx pgx.Tx, userID string) (any, error) {
	list, err := s.devices.List(ctx, tx, userID)
	if list == nil {
		list = []Device{}
	}
	return list, err
}
