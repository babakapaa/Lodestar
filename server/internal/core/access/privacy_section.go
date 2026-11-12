package access

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type PrivacySection struct {
	assignments AssignmentStore
}

func (PrivacySection) Name() string {
	return "roles"
}

func (s PrivacySection) Export(ctx context.Context, tx pgx.Tx, userID string) (any, error) {
	list, err := s.assignments.ByUser(ctx, tx, userID)
	if list == nil {
		list = []Assignment{}
	}
	return list, err
}
