package setup_test

import (
	"context"
	"slices"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/module"
)

type fakePilot struct {
	courses []module.PilotCourse
	members map[string][]string
}

func (f *fakePilot) PilotCourses(context.Context, pgx.Tx, string) ([]module.PilotCourse, error) {
	return f.courses, nil
}

func (f *fakePilot) InPilot(_ context.Context, _ pgx.Tx, userID string, courseIDs []string) (bool, error) {
	for _, id := range courseIDs {
		if slices.Contains(f.members[id], userID) {
			return true, nil
		}
	}
	return false, nil
}
