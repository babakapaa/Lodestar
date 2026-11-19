package demo

import (
	"context"
	"fmt"
	"time"

	"github.com/babakapaa/lodestar/server/internal/app"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/module"
)

func seedExtras(ctx context.Context, core app.Core, registry *module.Registry, tenantID string, ids map[string]string, courseID string, now time.Time) error {
	k := core.Kernel
	k.Help = registry.HelpSource(k)
	k.Briefs = registry.Briefs(k)
	k.Scores = registry.ScoreSources(k)
	k.TeamScopes = registry.TeamScopes(k)
	k.Teams = registry.TeamDirectory(k)
	k.Rosters = registry.Rosters(k)
	k.Lessons = registry.LessonSource(k)
	core.Kernel = k
	people := map[string]auth.Principal{}
	for subject, id := range ids {
		people[subject] = auth.Principal{TenantID: tenantID, UserID: id}
	}
	teacher := people[Teacher.Subject]
	if err := seedMaterials(ctx, core, teacher, courseID); err != nil {
		return fmt.Errorf("demo materials: %w", err)
	}
	if err := seedChats(ctx, core, people, courseID); err != nil {
		return fmt.Errorf("demo chats: %w", err)
	}
	if err := seedLeague(ctx, core, teacher, courseID, now); err != nil {
		return fmt.Errorf("demo league: %w", err)
	}
	return nil
}
