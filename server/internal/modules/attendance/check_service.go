package attendance

import (
	"context"
	"errors"
	"net/netip"
	"slices"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
)

func (s Service) Check(ctx context.Context, p auth.Principal, in CheckInput, from netip.Addr) (CheckResult, error) {
	in, err := in.normalized()
	if err != nil {
		return CheckResult{}, err
	}
	var out CheckResult
	err = s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		var err error
		out, err = s.check(ctx, tx, p, in, from)
		return err
	})
	if errors.Is(err, ErrWrongCode) {
		if failed := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
			return s.failures.Add(ctx, tx, p.UserID, s.now())
		}); failed != nil {
			return out, failed
		}
	}
	return out, err
}

func (s Service) check(ctx context.Context, tx pgx.Tx, p auth.Principal, in CheckInput, from netip.Addr) (CheckResult, error) {
	now := s.now()
	methods, err := s.methods(ctx, tx)
	if err != nil {
		return CheckResult{}, err
	}
	method := in.method()
	if method == ByQR && !methods.QR || method == ByCode && !methods.Code {
		return CheckResult{}, ErrMethodOff
	}
	if err := s.onCampus(ctx, tx, from); err != nil {
		return CheckResult{}, err
	}
	tries, err := s.failures.Count(ctx, tx, p.UserID, now.Add(-tryWindow))
	if err != nil {
		return CheckResult{}, err
	}
	if tries >= tryLimit {
		return CheckResult{}, ErrTooManyTries
	}
	r, err := s.target(ctx, tx, p, in, now)
	if err != nil {
		return CheckResult{}, err
	}
	if now.After(r.EndsAt) {
		return CheckResult{}, ErrSessionClosed
	}
	mark, known, err := s.marks.One(ctx, tx, r.ID, p.UserID)
	if err != nil {
		return CheckResult{}, err
	}
	if !known {
		l, err := s.lesson(ctx, tx, r.SeriesID, r.Date)
		if err != nil {
			return CheckResult{}, err
		}
		if !slices.Contains(l.Members, p.UserID) {
			return CheckResult{}, ErrNotInLesson
		}
	}
	out := CheckResult{SessionID: r.ID, CourseTitle: r.CourseTitle, StartsAt: r.StartsAt}
	if mark.Status != nil {
		out.Status = *mark.Status
		return out, nil
	}
	late, err := s.integer(ctx, tx, LateMinutes)
	if err != nil {
		return out, err
	}
	out.Status = Present
	if now.After(r.StartsAt.Add(time.Duration(late) * time.Minute)) {
		out.Status = Late
	}
	return out, s.marks.Put(ctx, tx, r.ID, p.UserID, &out.Status, &method, p.UserID, now)
}

func (s Service) target(ctx context.Context, tx pgx.Tx, p auth.Principal, in CheckInput, now time.Time) (sessionRow, error) {
	if in.Code == "" {
		r, err := s.sessions.ByID(ctx, tx, in.SessionID)
		if errors.Is(err, ErrSessionNotFound) {
			return r, ErrWrongCode
		}
		if err != nil {
			return r, err
		}
		if !r.open() {
			return r, ErrSessionClosed
		}
		if !qrValid(r, in.Token, now) {
			return r, ErrWrongCode
		}
		return r, nil
	}
	list, err := s.sessions.WithCode(ctx, tx, in.Code, now)
	if err != nil {
		return sessionRow{}, err
	}
	for _, r := range list {
		_, known, err := s.marks.One(ctx, tx, r.ID, p.UserID)
		if err != nil {
			return r, err
		}
		if known {
			return r, nil
		}
	}
	for _, r := range list {
		l, err := s.lesson(ctx, tx, r.SeriesID, r.Date)
		if err != nil {
			return r, err
		}
		if slices.Contains(l.Members, p.UserID) {
			return r, nil
		}
	}
	return sessionRow{}, ErrWrongCode
}

func (s Service) onCampus(ctx context.Context, tx pgx.Tx, from netip.Addr) error {
	networks, err := s.networks(ctx, tx)
	if err != nil || len(networks) == 0 {
		return err
	}
	for _, n := range networks {
		if from.IsValid() && n.Contains(from) {
			return nil
		}
	}
	return ErrOffCampus
}
