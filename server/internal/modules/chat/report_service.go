package chat

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/module"
)

type ReportService struct {
	k        module.Kernel
	chats    ChatStore
	messages MessageStore
	reports  ReportStore
	people   identity.People
}

func NewReportService(k module.Kernel) ReportService {
	sealer := Sealer{box: k.Secrets}
	return ReportService{k: k, chats: ChatStore{sealer: sealer}, messages: MessageStore{sealer: sealer}, reports: ReportStore{sealer: sealer}}
}

func (s ReportService) Create(ctx context.Context, p auth.Principal, in ReportInput) (Report, error) {
	in, err := in.normalized()
	if err != nil {
		return Report{}, err
	}
	var r Report
	err = s.k.Tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		c, err := s.chats.ForMember(ctx, tx, p.UserID, in.ChatID)
		if err != nil {
			return err
		}
		body := reportBody{ChatTitle: c.Title, Reason: in.Reason}
		var authors []string
		for _, id := range in.MessageIDs {
			m, err := s.messages.ByID(ctx, tx, id)
			if errors.Is(err, ErrMessageNotFound) || (err == nil && (m.ChatID != c.ID || m.Deleted)) {
				return ErrInvalidReport
			}
			if err != nil {
				return err
			}
			body.Messages = append(body.Messages, ReportedMessage{Seq: m.Seq, AuthorID: m.AuthorID, Text: m.Text, CreatedAt: m.CreatedAt})
			if m.AuthorID != nil {
				authors = append(authors, *m.AuthorID)
			}
		}
		people, err := s.people.ByIDs(ctx, tx, authors)
		if err != nil {
			return err
		}
		for i, m := range body.Messages {
			for _, person := range people {
				if m.AuthorID != nil && *m.AuthorID == person.UserID {
					body.Messages[i].AuthorName = person.DisplayName
				}
			}
		}
		r, err = s.reports.Insert(ctx, tx, c.ID, p.UserID, body)
		return err
	})
	return r, err
}

func (s ReportService) Open(ctx context.Context, p auth.Principal) ([]Report, error) {
	var list []Report
	err := s.reviewer(ctx, p, func(tx pgx.Tx) error {
		var err error
		list, err = s.reports.Open(ctx, tx)
		return err
	})
	return list, err
}

func (s ReportService) Resolve(ctx context.Context, p auth.Principal, id string) (Report, error) {
	if !validID(id) {
		return Report{}, ErrReportNotFound
	}
	var r Report
	err := s.reviewer(ctx, p, func(tx pgx.Tx) error {
		var err error
		r, err = s.reports.Resolve(ctx, tx, id, p.UserID)
		return err
	})
	return r, err
}

func (s ReportService) reviewer(ctx context.Context, p auth.Principal, fn func(tx pgx.Tx) error) error {
	return s.k.Tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		allowed, err := s.k.Access.AllowedInTenant(ctx, tx, p, string(ReviewReports))
		if err != nil {
			return err
		}
		if !allowed {
			return access.ErrForbidden
		}
		return fn(tx)
	})
}
