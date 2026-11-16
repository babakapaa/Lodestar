package demo

import (
	"context"
	"fmt"
	"time"

	"github.com/babakapaa/lodestar/server/internal/app"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/structure"
	"github.com/babakapaa/lodestar/server/internal/modules/study"
)

const (
	CourseCode = "BIO-101"
	GroupName  = "ДЕМО-21"
)

func seedStudy(ctx context.Context, core app.Core, tenantID string, ids map[string]string, now time.Time) (string, bool, error) {
	admin := auth.Principal{TenantID: tenantID, UserID: ids[Admin.Subject]}
	teacher := auth.Principal{TenantID: tenantID, UserID: ids[Teacher.Subject]}
	tx := core.Kernel.Tenants
	courses := study.NewService(tx, core.Access)
	existing, err := courses.Catalog(ctx, admin, nil, true)
	if err != nil {
		return "", false, err
	}
	units := structure.NewService(tx, core.Access)
	for _, c := range existing {
		if c.Code != nil && *c.Code == CourseCode {
			return c.ID, false, joinGroup(ctx, units, admin, ids[Peer.Subject])
		}
	}
	group, err := units.Create(ctx, admin, structure.UnitInput{Kind: "group", Name: GroupName})
	if err != nil {
		return "", false, err
	}
	for _, student := range []string{ids[Student.Subject], ids[Peer.Subject]} {
		if err := units.AddMember(ctx, admin, group.ID, student); err != nil {
			return "", false, err
		}
	}
	code := CourseCode
	course, err := courses.Create(ctx, admin, study.CourseInput{
		Title:       "Основы биологии",
		Code:        &code,
		Description: "Клетка, ткани и обмен веществ. Курс для демонстрации Lodestar.",
	})
	if err != nil {
		return "", false, err
	}
	if err := courses.Change(ctx, admin, study.ListTeachers, course.ID, ids[Teacher.Subject], true); err != nil {
		return "", false, err
	}
	if err := courses.Change(ctx, admin, study.ListUnits, course.ID, group.ID, true); err != nil {
		return "", false, err
	}
	if err := seedAssignments(ctx, core, teacher, course.ID, now); err != nil {
		return "", false, fmt.Errorf("demo assignments: %w", err)
	}
	schedule := study.NewScheduleService(core.Kernel)
	today := study.DayOf(now)
	for _, lesson := range []study.SeriesInput{
		{Kind: "lecture", Room: "А-204", Weekday: 1, StartTime: study.Clock(9 * 60)},
		{Kind: "practice", Room: "Б-112", Weekday: 3, StartTime: study.Clock(11 * 60)},
	} {
		lesson.UnitID = &group.ID
		lesson.TeacherID = &teacher.UserID
		lesson.StartsOn = today.AddDays(-7)
		lesson.EndsOn = today.AddDays(120)
		lesson.Minutes = 90
		if _, err := schedule.Create(ctx, admin, course.ID, lesson); err != nil {
			return "", false, err
		}
	}
	return course.ID, true, nil
}

func joinGroup(ctx context.Context, units structure.Service, admin auth.Principal, userID string) error {
	groups, err := units.Children(ctx, admin, nil, false)
	if err != nil {
		return err
	}
	for _, g := range groups {
		if g.Name == GroupName {
			return units.AddMember(ctx, admin, g.ID, userID)
		}
	}
	return nil
}
