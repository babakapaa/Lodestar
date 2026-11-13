package setup_test

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/module"
)

type sentNotifications struct {
	list []module.Notification
}

func (s *sentNotifications) Notify(_ context.Context, _ pgx.Tx, _ string, n module.Notification) error {
	s.list = append(s.list, n)
	return nil
}

func (s *sentNotifications) take() []module.Notification {
	out := s.list
	s.list = nil
	return out
}
