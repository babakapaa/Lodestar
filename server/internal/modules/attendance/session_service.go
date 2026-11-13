package attendance

import (
	"context"
	"crypto/rand"
	"slices"
	"sort"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
)

func (s Service) Open(ctx context.Context, p auth.Principal, seriesID, date string) (Session, error) {
	var out Session
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		l, err := s.lesson(ctx, tx, seriesID, date)
		if err != nil {
			return err
		}
		if err := s.lead(ctx, tx, p, l.Leads); err != nil {
			return err
		}
		if l.Cancelled {
			return ErrLessonCancelled
		}
		secret := make([]byte, 32)
		if _, err := rand.Read(secret); err != nil {
			return err
		}
		code, err := s.freshCode(ctx, tx)
		if err != nil {
			return err
		}
		r, err := s.sessions.Open(ctx, tx, l, secret, code, s.now().Add(codeLife), p.UserID)
		if err != nil {
			return err
		}
		if err := s.marks.Expect(ctx, tx, r.ID, l.Members); err != nil {
			return err
		}
		if err := s.record(ctx, tx, p, ActionOpened, r.ID, map[string]any{"series": seriesID, "date": date}); err != nil {
			return err
		}
		out, err = s.view(ctx, tx, r)
		return err
	})
	return out, err
}

func (s Service) Session(ctx context.Context, p auth.Principal, id string) (Session, error) {
	var out Session
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		r, err := s.sessions.ByID(ctx, tx, id)
		if err != nil {
			return err
		}
		if _, err := s.leadOf(ctx, tx, p, r); err != nil {
			return err
		}
		out, err = s.view(ctx, tx, r)
		return err
	})
	return out, err
}

func (s Service) QR(ctx context.Context, p auth.Principal, id string) (QRToken, error) {
	var out QRToken
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		r, err := s.openSession(ctx, tx, p, id)
		if err != nil {
			return err
		}
		out = qrFor(r, s.now())
		return nil
	})
	return out, err
}

func (s Service) RefreshCode(ctx context.Context, p auth.Principal, id string) (Session, error) {
	var out Session
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		r, err := s.openSession(ctx, tx, p, id)
		if err != nil {
			return err
		}
		if r.Code, err = s.freshCode(ctx, tx); err != nil {
			return err
		}
		r.CodeUntil = s.now().Add(codeLife)
		if err := s.sessions.SetCode(ctx, tx, r.ID, r.Code, r.CodeUntil); err != nil {
			return err
		}
		out, err = s.view(ctx, tx, r)
		return err
	})
	return out, err
}

func (s Service) Mark(ctx context.Context, p auth.Principal, id, userID string, in MarkInput) (Session, error) {
	in, err := in.normalized()
	if err != nil {
		return Session{}, err
	}
	var out Session
	err = s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		r, err := s.sessions.ByID(ctx, tx, id)
		if err != nil {
			return err
		}
		l, err := s.leadOf(ctx, tx, p, r)
		if err != nil {
			return err
		}
		_, known, err := s.marks.One(ctx, tx, id, userID)
		if err != nil {
			return err
		}
		if !known && !slices.Contains(l.Members, userID) {
			return ErrNotInLesson
		}
		method := ByList
		if !r.open() {
			method = ByManual
		}
		var how *MarkMethod
		if in.Status != nil {
			how = &method
		}
		if err := s.marks.Put(ctx, tx, id, userID, in.Status, how, p.UserID, s.now()); err != nil {
			return err
		}
		if !r.open() {
			if err := s.record(ctx, tx, p, ActionCorrect, id, map[string]any{"user": userID, "status": in.Status}); err != nil {
				return err
			}
		}
		out, err = s.view(ctx, tx, r)
		return err
	})
	return out, err
}

func (s Service) Close(ctx context.Context, p auth.Principal, id string) (Session, error) {
	var out Session
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		r, err := s.openSession(ctx, tx, p, id)
		if err != nil {
			return err
		}
		now := s.now()
		if err := s.marks.FillAbsent(ctx, tx, id, p.UserID, now); err != nil {
			return err
		}
		if err := s.sessions.Close(ctx, tx, id, now); err != nil {
			return err
		}
		if err := s.record(ctx, tx, p, ActionClosed, id, nil); err != nil {
			return err
		}
		r.ClosedAt = &now
		out, err = s.view(ctx, tx, r)
		return err
	})
	return out, err
}

func (s Service) openSession(ctx context.Context, tx pgx.Tx, p auth.Principal, id string) (sessionRow, error) {
	r, err := s.sessions.ByID(ctx, tx, id)
	if err != nil {
		return r, err
	}
	if _, err := s.leadOf(ctx, tx, p, r); err != nil {
		return r, err
	}
	if !r.open() {
		return r, ErrSessionClosed
	}
	return r, nil
}

func (s Service) view(ctx context.Context, tx pgx.Tx, r sessionRow) (Session, error) {
	methods, err := s.methods(ctx, tx)
	if err != nil {
		return Session{}, err
	}
	out := Session{
		ID: r.ID, SeriesID: r.SeriesID, Date: r.Date, CourseID: r.CourseID, CourseTitle: r.CourseTitle, Room: r.Room,
		StartsAt: r.StartsAt, EndsAt: r.EndsAt, Open: r.open(), Methods: methods, Marks: []Mark{},
	}
	if out.Open && methods.Code {
		out.Code, out.CodeUntil = r.Code, &r.CodeUntil
	}
	marks, err := s.marks.OfSessions(ctx, tx, []string{r.ID})
	if err != nil {
		return out, err
	}
	ids := make([]string, len(marks))
	for i, m := range marks {
		ids[i] = m.UserID
	}
	names, err := s.names(ctx, tx, ids)
	if err != nil {
		return out, err
	}
	for _, m := range marks {
		out.Counts.add(m.Status)
		out.Marks = append(out.Marks, Mark{UserID: m.UserID, Name: names[m.UserID], Status: m.Status, Method: m.Method, MarkedAt: m.MarkedAt})
	}
	sort.SliceStable(out.Marks, func(i, j int) bool { return out.Marks[i].Name < out.Marks[j].Name })
	return out, nil
}
