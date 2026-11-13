package navigator

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/module"
)

const scopeCourse = "course"

type MaterialService struct {
	k         module.Kernel
	materials MaterialStore
}

func NewMaterialService(k module.Kernel) MaterialService {
	return MaterialService{k: k}
}

func (s MaterialService) List(ctx context.Context, p auth.Principal, courseID string) ([]Material, error) {
	var list []Material
	err := s.k.Tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if _, err := s.k.Briefs.Brief(ctx, tx, p, scopeCourse, courseID); err != nil {
			return err
		}
		var err error
		list, err = s.materials.List(ctx, tx, s.model(), courseID)
		return err
	})
	return list, err
}

func (s MaterialService) Add(ctx context.Context, p auth.Principal, courseID string, in MaterialInput) (Material, error) {
	in, err := in.normalized()
	if err != nil {
		return Material{}, err
	}
	var m Material
	err = s.k.Tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if err := s.taught(ctx, tx, p, courseID); err != nil {
			return err
		}
		m, err = s.materials.Insert(ctx, tx, s.model(), courseID, p.UserID, in, chunk(in.Text))
		return err
	})
	return m, err
}

func (s MaterialService) Delete(ctx context.Context, p auth.Principal, id string) error {
	return s.k.Tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		m, err := s.materials.ByID(ctx, tx, s.model(), id)
		if err != nil {
			return err
		}
		if err := s.taught(ctx, tx, p, m.CourseID); err != nil {
			return err
		}
		return s.materials.Delete(ctx, tx, id)
	})
}

func (s MaterialService) taught(ctx context.Context, tx pgx.Tx, p auth.Principal, courseID string) error {
	b, err := s.k.Briefs.Brief(ctx, tx, p, scopeCourse, courseID)
	if err != nil {
		return err
	}
	if !b.Teaches {
		return access.ErrForbidden
	}
	return nil
}

func (s MaterialService) model() string {
	if s.k.Embedder == nil {
		return ""
	}
	return s.k.Embedder.Name()
}
