package privacy

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

type Service struct {
	tenants   database.TenantTx
	settings  Settings
	deletions DeletionStore
}

func NewService(tenants database.TenantTx, settings Settings) Service {
	return Service{tenants: tenants, settings: settings}
}

func (s Service) RequestDeletion(ctx context.Context, p auth.Principal) (Deletion, error) {
	var d Deletion
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		var allowed bool
		if err := s.read(ctx, tx, SelfDeletionSetting, &allowed); err != nil {
			return err
		}
		if !allowed {
			return ErrSelfDeletionOff
		}
		var days int
		if err := s.read(ctx, tx, GraceDaysSetting, &days); err != nil {
			return err
		}
		var err error
		d, err = s.deletions.Put(ctx, tx, p.UserID, time.Duration(days)*24*time.Hour)
		return err
	})
	return d, err
}

func (s Service) Deletion(ctx context.Context, p auth.Principal) (Deletion, error) {
	var d Deletion
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		var err error
		d, err = s.deletions.Get(ctx, tx, p.UserID)
		return err
	})
	return d, err
}

func (s Service) CancelDeletion(ctx context.Context, p auth.Principal) error {
	return s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		found, err := s.deletions.Delete(ctx, tx, p.UserID)
		if err != nil {
			return err
		}
		if !found {
			return ErrNoRequest
		}
		return nil
	})
}

func (s Service) read(ctx context.Context, tx pgx.Tx, key string, v any) error {
	raw, err := s.settings.Value(ctx, tx, key)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(raw, v); err != nil {
		return fmt.Errorf("read %s: %w", key, err)
	}
	return nil
}
