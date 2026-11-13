package sandbox

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/audit"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/identity"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

type Service struct {
	tenants   database.TenantTx
	access    access.Service
	settings  module.Settings
	notifier  module.Notifier
	briefs    module.Briefs
	teams     module.TeamDirectory
	grants    module.GPUSources
	classes   ClassStore
	terms     TermStore
	ledger    LedgerStore
	topups    TopupStore
	templates TemplateStore
	nodes     NodeStore
	runs      RunStore
	people    identity.People
	log       audit.Log
}

func NewService(k module.Kernel) Service {
	return Service{
		tenants: k.Tenants, access: k.Access, settings: k.Settings, notifier: k.Notifier,
		briefs: k.Briefs, teams: k.Teams, grants: k.GPU,
	}
}

func (s Service) manage(ctx context.Context, tx pgx.Tx, p auth.Principal) error {
	return s.access.Require(ctx, tx, p, access.TenantScope(p.TenantID), Manage)
}

func (s Service) current(ctx context.Context, tx pgx.Tx) (Term, error) {
	t, ok, err := s.terms.Current(ctx, tx)
	if err != nil {
		return t, err
	}
	if !ok {
		return t, ErrNoTerm
	}
	return t, nil
}

func (s Service) name(ctx context.Context, tx pgx.Tx, userID string) (string, error) {
	people, err := s.people.ByIDs(ctx, tx, []string{userID})
	if err != nil || len(people) == 0 {
		return "", err
	}
	if people[0].DisplayName != "" {
		return people[0].DisplayName, nil
	}
	return people[0].Email, nil
}
