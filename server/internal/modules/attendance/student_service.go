package attendance

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
)

func (s Service) OpenLessons(ctx context.Context, p auth.Principal) ([]OpenLesson, error) {
	out := []OpenLesson{}
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		methods, err := s.methods(ctx, tx)
		if err != nil {
			return err
		}
		rows, err := s.sessions.OpenFor(ctx, tx, p.UserID, s.now())
		if err != nil {
			return err
		}
		for _, o := range rows {
			out = append(out, OpenLesson{
				SessionID: o.session.ID, CourseTitle: o.session.CourseTitle, Room: o.session.Room,
				StartsAt: o.session.StartsAt, EndsAt: o.session.EndsAt, Status: o.status, Methods: methods,
			})
		}
		return nil
	})
	return out, err
}

func (s Service) Mine(ctx context.Context, p auth.Principal) ([]MyMark, error) {
	var out []MyMark
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		var err error
		out, err = s.marks.Mine(ctx, tx, p.UserID)
		return err
	})
	if out == nil {
		out = []MyMark{}
	}
	return out, err
}
