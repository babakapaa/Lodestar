package lti_test

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/module"
)

type fakeLinks struct {
	owner *string
}

func (fakeLinks) Kind() string {
	return "study.course"
}

func (f fakeLinks) Linkable(_ context.Context, _ pgx.Tx, p auth.Principal) ([]module.LinkTarget, error) {
	if p.UserID != *f.owner {
		return []module.LinkTarget{}, nil
	}
	return []module.LinkTarget{f.target("bio")}, nil
}

func (f fakeLinks) Describe(_ context.Context, _ pgx.Tx, id string) (module.LinkTarget, error) {
	if id != "bio" {
		return module.LinkTarget{}, module.ErrLinkTargetNotFound
	}
	return f.target(id), nil
}

func (f fakeLinks) CanLink(_ context.Context, _ pgx.Tx, p auth.Principal, id string) error {
	if id != "bio" {
		return module.ErrLinkTargetNotFound
	}
	if p.UserID != *f.owner {
		return access.ErrForbidden
	}
	return nil
}

func (fakeLinks) target(id string) module.LinkTarget {
	return module.LinkTarget{Kind: "study.course", ID: id, Title: "Биология", Path: "/study/courses/" + id}
}
