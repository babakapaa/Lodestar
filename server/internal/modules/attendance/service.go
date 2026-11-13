package attendance

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/audit"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

const (
	codeLife   = 2 * time.Minute
	tryWindow  = 10 * time.Minute
	tryLimit   = 5
	codeChoice = 20
)

type Service struct {
	tenants  database.TenantTx
	access   access.Service
	settings module.Settings
	lessons  module.LessonSource
	sessions SessionStore
	marks    MarkStore
	failures FailureStore
	people   identity.People
	log      audit.Log
	now      func() time.Time
}

func NewService(k module.Kernel) Service {
	return Service{tenants: k.Tenants, access: k.Access, settings: k.Settings, lessons: k.Lessons, now: time.Now}
}

func (s Service) WithClock(now func() time.Time) Service {
	s.now = now
	return s
}

func (s Service) lesson(ctx context.Context, tx pgx.Tx, seriesID, date string) (module.LessonRef, error) {
	if s.lessons == nil {
		return module.LessonRef{}, ErrLessonNotFound
	}
	l, err := s.lessons.Lesson(ctx, tx, seriesID, date)
	if errors.Is(err, module.ErrLessonNotFound) {
		return l, ErrLessonNotFound
	}
	return l, err
}

func (s Service) lead(ctx context.Context, tx pgx.Tx, p auth.Principal, leads []string) error {
	if slices.Contains(leads, p.UserID) {
		return nil
	}
	return s.access.Require(ctx, tx, p, access.TenantScope(p.TenantID), Manage)
}

func (s Service) leadOf(ctx context.Context, tx pgx.Tx, p auth.Principal, r sessionRow) (module.LessonRef, error) {
	l, err := s.lesson(ctx, tx, r.SeriesID, r.Date)
	if err != nil {
		return l, err
	}
	if err := s.lead(ctx, tx, p, l.Leads); err != nil {
		if errors.Is(err, access.ErrForbidden) {
			return l, ErrSessionNotFound
		}
		return l, err
	}
	return l, nil
}

func (s Service) freshCode(ctx context.Context, tx pgx.Tx) (string, error) {
	now := s.now()
	for range codeChoice {
		n, err := rand.Int(rand.Reader, big.NewInt(10000))
		if err != nil {
			return "", err
		}
		code := fmt.Sprintf("%04d", n.Int64())
		taken, err := s.sessions.CodeTaken(ctx, tx, code, now)
		if err != nil {
			return "", err
		}
		if !taken {
			return code, nil
		}
	}
	return "", errors.New("no free attendance code")
}

func (s Service) names(ctx context.Context, tx pgx.Tx, ids []string) (map[string]string, error) {
	people, err := s.people.ByIDs(ctx, tx, ids)
	if err != nil {
		return nil, err
	}
	out := make(map[string]string, len(people))
	for _, p := range people {
		out[p.UserID] = p.DisplayName
		if p.DisplayName == "" {
			out[p.UserID] = p.Email
		}
	}
	return out, nil
}

func (s Service) record(ctx context.Context, tx pgx.Tx, p auth.Principal, action, sessionID string, details any) error {
	return s.log.Record(ctx, tx, p.UserID, audit.Change{Action: action, TargetKind: "attendance.session", TargetID: sessionID, Details: details})
}
