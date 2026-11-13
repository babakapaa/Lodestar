package sandbox

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/audit"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
)

func (s Service) Classes(ctx context.Context, p auth.Principal) ([]GPUClass, error) {
	var list []GPUClass
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		var err error
		list, err = s.classes.All(ctx, tx)
		return err
	})
	return list, err
}

func (s Service) CreateClass(ctx context.Context, p auth.Principal, in GPUClassInput) (GPUClass, error) {
	return s.saveClass(ctx, p, "", in)
}

func (s Service) UpdateClass(ctx context.Context, p auth.Principal, id string, in GPUClassInput) (GPUClass, error) {
	return s.saveClass(ctx, p, id, in)
}

func (s Service) saveClass(ctx context.Context, p auth.Principal, id string, in GPUClassInput) (GPUClass, error) {
	in, err := in.normalized()
	if err != nil {
		return GPUClass{}, err
	}
	var c GPUClass
	err = s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if err := s.manage(ctx, tx, p); err != nil {
			return err
		}
		if id == "" {
			c, err = s.classes.Insert(ctx, tx, in)
		} else {
			c, err = s.classes.Update(ctx, tx, id, in)
		}
		if err != nil {
			return err
		}
		return s.log.Record(ctx, tx, p.UserID, audit.Change{Action: ActionClass, TargetKind: "sandbox.class", TargetID: c.ID, Details: in})
	})
	return c, err
}
