package demo

import (
	"context"
	"time"

	"github.com/babakapaa/lodestar/server/internal/app"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/modules/study"
)

func seedAssignments(ctx context.Context, core app.Core, teacher auth.Principal, courseID string, now time.Time) error {
	assignments := study.NewAssignmentService(core.Kernel.Tenants, core.Access, core.Kernel.Notifier)
	essayDue := now.Add(5 * 24 * time.Hour).Truncate(time.Hour)
	essayPoints := 10.0
	defense := study.DefenseShort
	_, err := assignments.Create(ctx, teacher, courseID, study.AssignmentInput{
		Title:       "Эссе о строении клетки",
		Description: "Опишите строение животной клетки и роль каждого органоида. Объём до двух страниц. Перед сдачей соберите карту поиска: откуда вы брали сведения и почему им можно верить.",
		DueAt:       &essayDue,
		Grading:     study.GradingPoints,
		MaxPoints:   &essayPoints,
		Weights:     &study.GradeWeights{Defense: 50, Map: 30, Work: 20},
		SearchMap:   &study.MapRules{MinSources: 3, MinScholarly: 1},
		Defense:     &defense,
		Published:   true,
	})
	if err != nil {
		return err
	}
	projectDue := now.Add(12 * 24 * time.Hour).Truncate(time.Hour)
	_, err = assignments.Create(ctx, teacher, courseID, study.AssignmentInput{
		Title:       "Групповой проект: модель клетки",
		Description: "Соберите в команде объёмную или цифровую модель клетки и подпишите органоиды. Команду из 2–3 человек соберите сами.",
		DueAt:       &projectDue,
		Grading:     study.GradingPass,
		Team:        &study.TeamRules{MinSize: 2, MaxSize: 3, SelfFormed: true},
		Published:   true,
	})
	return err
}
