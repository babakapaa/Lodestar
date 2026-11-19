package navigator

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/module"
)

const (
	privacyFloor = 5
	recentDays   = 14
	defaultWeeks = 8
	longestWeeks = 52
)

type TopicService struct {
	k      module.Kernel
	topics TopicStore
	now    func() time.Time
}

func NewTopicService(k module.Kernel) TopicService {
	return TopicService{k: k, now: time.Now}
}

func (s TopicService) Course(ctx context.Context, p auth.Principal, courseID string, weeks int) (CourseTopics, error) {
	if weeks <= 0 {
		weeks = defaultWeeks
	}
	weeks = min(weeks, longestWeeks)
	out := CourseTopics{Threshold: privacyFloor, Weeks: weeks, Topics: []TopicCount{}, Intents: []IntentCount{}, Trend: []WeekCount{}}
	err := s.k.Tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		b, err := s.k.Briefs.Brief(ctx, tx, p, scopeCourse, courseID)
		if err != nil {
			return err
		}
		if !b.Teaches {
			return access.ErrForbidden
		}
		today := s.now().UTC().Truncate(24 * time.Hour)
		monday := today.AddDate(0, 0, -((int(today.Weekday()) + 6) % 7))
		from := monday.AddDate(0, 0, -7*(weeks-1))
		total, err := s.topics.Students(ctx, tx, courseID, from)
		if err != nil {
			return err
		}
		if total < privacyFloor {
			return nil
		}
		out.Students = &total
		if out.Topics, err = s.topics.Topics(ctx, tx, courseID, from, today.AddDate(0, 0, -recentDays), privacyFloor); err != nil {
			return err
		}
		if out.Intents, err = s.topics.Intents(ctx, tx, courseID, from, privacyFloor); err != nil {
			return err
		}
		weekly, err := s.topics.Weekly(ctx, tx, courseID, from)
		if err != nil {
			return err
		}
		for w := from; !w.After(monday); w = w.AddDate(0, 0, 7) {
			key := w.Format(time.DateOnly)
			c := WeekCount{Week: key}
			if n := weekly[key]; n >= privacyFloor {
				c.Students = &n
			}
			out.Trend = append(out.Trend, c)
		}
		for i := range out.Topics {
			if out.Topics[i].Recent < privacyFloor {
				out.Topics[i].Recent = 0
			}
		}
		for i, t := range out.Topics {
			if t.Recent >= privacyFloor && (out.Suggestion == nil || t.Recent > out.Suggestion.Recent) {
				out.Suggestion = &out.Topics[i]
			}
		}
		return nil
	})
	return out, err
}
