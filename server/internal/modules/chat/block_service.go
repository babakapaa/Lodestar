package chat

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/module"
)

type BlockService struct {
	k      module.Kernel
	blocks BlockStore
	people identity.People
}

func NewBlockService(k module.Kernel) BlockService {
	return BlockService{k: k}
}

func (s BlockService) List(ctx context.Context, p auth.Principal) ([]identity.Person, error) {
	var list []identity.Person
	err := s.k.Tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		ids, err := s.blocks.List(ctx, tx, p.UserID)
		if err != nil {
			return err
		}
		list, err = s.people.ByIDs(ctx, tx, ids)
		return err
	})
	return list, err
}

func (s BlockService) Block(ctx context.Context, p auth.Principal, in DirectInput) error {
	if !validID(in.UserID) || in.UserID == p.UserID {
		return ErrInvalidChat
	}
	return s.k.Tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		return s.blocks.Add(ctx, tx, p.UserID, in.UserID)
	})
}

func (s BlockService) Unblock(ctx context.Context, p auth.Principal, userID string) error {
	if !validID(userID) {
		return ErrInvalidChat
	}
	return s.k.Tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		return s.blocks.Remove(ctx, tx, p.UserID, userID)
	})
}
