package lti_test

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/module"
)

type fakeLinks struct {
	mu       sync.Mutex
	owner    string
	manager  string
	leads    []string
	learners []string
	works    []module.ExternalWork
	updated  []string
}

func (*fakeLinks) Kind() string {
	return "study.course"
}

func (f *fakeLinks) Linkable(_ context.Context, _ pgx.Tx, p auth.Principal) ([]module.LinkTarget, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if p.UserID != f.owner {
		return []module.LinkTarget{}, nil
	}
	return []module.LinkTarget{f.target("bio")}, nil
}

func (f *fakeLinks) Describe(_ context.Context, _ pgx.Tx, id string) (module.LinkTarget, error) {
	if id != "bio" {
		return module.LinkTarget{}, module.ErrLinkTargetNotFound
	}
	return f.target(id), nil
}

func (f *fakeLinks) CanLink(_ context.Context, _ pgx.Tx, p auth.Principal, id string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if id != "bio" {
		return module.ErrLinkTargetNotFound
	}
	if p.UserID != f.owner {
		return access.ErrForbidden
	}
	return nil
}

func (f *fakeLinks) CanManage(_ context.Context, _ pgx.Tx, p auth.Principal, id string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if id != "bio" {
		return module.ErrLinkTargetNotFound
	}
	if p.UserID != f.manager {
		return access.ErrForbidden
	}
	return nil
}

func (f *fakeLinks) Enroll(_ context.Context, _ pgx.Tx, _, _ string, leads, members []string) (int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.leads, f.learners = leads, members
	return len(leads) + len(members), nil
}

func (f *fakeLinks) PutWork(_ context.Context, _ pgx.Tx, _, _, workID string, w module.ExternalWork) (module.WorkRef, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if workID != "" {
		f.updated = append(f.updated, workID)
		return module.WorkRef{Kind: "study.assignment", ID: workID}, nil
	}
	f.works = append(f.works, w)
	return module.WorkRef{Kind: "study.assignment", ID: "w-" + w.Title}, nil
}

func (f *fakeLinks) target(id string) module.LinkTarget {
	return module.LinkTarget{Kind: "study.course", ID: id, Title: "Биология", Path: "/study/courses/" + id}
}

func (f *fakeLinks) ScoredWorks(context.Context, pgx.Tx, string) ([]module.ScoredWork, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	w := module.ScoredWork{Kind: "study.assignment", ID: "essay", Title: "Эссе о клетке", Max: 10}
	for _, id := range f.learners {
		w.Results = append(w.Results, module.WorkResult{UserID: id, Score: 8, Max: 10})
	}
	return []module.ScoredWork{w}, nil
}

func (*fakeLinks) All(context.Context, pgx.Tx) ([]module.Roster, error) {
	return nil, nil
}

func (*fakeLinks) One(context.Context, pgx.Tx, string) (module.Roster, error) {
	return module.Roster{}, module.ErrRosterNotFound
}

func (f *fakeLinks) people(owner, manager string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.owner, f.manager = owner, manager
}
