package delivery

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

type Service struct {
	tenants database.TenantTx
	devices DeviceStore
}

func NewService(tenants database.TenantTx) Service {
	return Service{tenants: tenants}
}

func (s Service) Register(ctx context.Context, p auth.Principal, in DeviceInput) (Device, error) {
	in, err := in.normalized()
	if err != nil {
		return Device{}, err
	}
	var d Device
	err = s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		var err error
		d, err = s.devices.Upsert(ctx, tx, p.UserID, in)
		return err
	})
	return d, err
}

func (s Service) Devices(ctx context.Context, p auth.Principal) ([]Device, error) {
	var list []Device
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		var err error
		list, err = s.devices.List(ctx, tx, p.UserID)
		return err
	})
	return list, err
}

func (s Service) Remove(ctx context.Context, p auth.Principal, id string) error {
	return s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		found, err := s.devices.Delete(ctx, tx, p.UserID, id)
		if err != nil {
			return err
		}
		if !found {
			return ErrDeviceNotFound
		}
		return nil
	})
}
