package demo

import (
	"context"
	"time"

	"github.com/babakapaa/lodestar/server/internal/app"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/modules/league"
)

const courseArena = "study.course"

func seedLeague(ctx context.Context, core app.Core, teacher auth.Principal, courseID string, now time.Time) error {
	leagues := league.NewService(core.Kernel)
	view, err := leagues.View(ctx, teacher, courseArena, courseID)
	if err != nil {
		return err
	}
	due := now.Add(14 * 24 * time.Hour).Truncate(time.Hour)
	_, err = leagues.CreateTask(ctx, teacher, view.ID, league.TaskInput{
		Title:      "Плакат «Как устроена клетка» для школьников",
		About:      "Подготовьте понятный плакат для учеников 7 класса. Лучшие работы покажем на дне открытых дверей.",
		Partner:    "Школа-лицей № 1",
		Points:     200,
		MinSize:    2,
		MaxSize:    3,
		SelfFormed: true,
		DueAt:      &due,
	})
	return err
}
