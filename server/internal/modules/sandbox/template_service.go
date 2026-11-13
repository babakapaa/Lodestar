package sandbox

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/audit"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
)

func (s Service) Templates(ctx context.Context, p auth.Principal) ([]Template, error) {
	var list []Template
	err := s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		var err error
		list, err = s.templates.All(ctx, tx)
		return err
	})
	return list, err
}

func (s Service) CreateTemplate(ctx context.Context, p auth.Principal, in TemplateInput) (Template, error) {
	return s.saveTemplate(ctx, p, "", in)
}

func (s Service) UpdateTemplate(ctx context.Context, p auth.Principal, id string, in TemplateInput) (Template, error) {
	return s.saveTemplate(ctx, p, id, in)
}

func (s Service) saveTemplate(ctx context.Context, p auth.Principal, id string, in TemplateInput) (Template, error) {
	in, err := in.normalized()
	if err != nil {
		return Template{}, err
	}
	var t Template
	err = s.tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if err := s.manage(ctx, tx, p); err != nil {
			return err
		}
		if id == "" {
			t, err = s.templates.Insert(ctx, tx, in)
		} else {
			t, err = s.templates.Update(ctx, tx, id, in)
		}
		if err != nil {
			return err
		}
		return s.log.Record(ctx, tx, p.UserID, audit.Change{Action: ActionTemplate, TargetKind: "sandbox.template", TargetID: t.ID, Details: in})
	})
	return t, err
}
