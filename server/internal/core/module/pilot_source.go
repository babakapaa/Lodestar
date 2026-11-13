package module

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type PilotSource interface {
	PilotCourses(ctx context.Context, tx pgx.Tx, tenantID string) ([]PilotCourse, error)
	InPilot(ctx context.Context, tx pgx.Tx, userID string, courseIDs []string) (bool, error)
}
