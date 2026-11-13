package lti

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
)

const (
	lodestarTag  = "lodestar"
	maxSyncError = 500
)

type SyncService struct {
	tx         database.TenantTx
	values     Values
	links      module.LinkSources
	scores     module.ScoreSources
	client     *PlatformClient
	accounts   Accounts
	keys       KeyStore
	platforms  PlatformStore
	contexts   ContextStore
	members    MemberStore
	identities IdentityStore
	lineItems  LineItemStore
	sent       ScoreStore
}

func (s SyncService) Sync(ctx context.Context, tenantID, contextID string) (SyncReport, error) {
	plan, err := s.plan(ctx, tenantID, contextID)
	if err != nil {
		return SyncReport{}, err
	}
	report, syncErr := s.run(ctx, plan)
	var text *string
	if syncErr != nil {
		msg := syncErr.Error()
		if len(msg) > maxSyncError {
			msg = msg[:maxSyncError]
		}
		text = &msg
	}
	err = s.tx.Run(ctx, tenantID, func(tx pgx.Tx) error {
		return s.contexts.Synced(ctx, tx, contextID, text)
	})
	if syncErr != nil {
		return report, syncErr
	}
	return report, err
}

func (s SyncService) plan(ctx context.Context, tenantID, contextID string) (syncPlan, error) {
	plan := syncPlan{tenantID: tenantID}
	err := s.tx.Run(ctx, tenantID, func(tx pgx.Tx) error {
		var err error
		if plan.context, err = s.contexts.ByID(ctx, tx, contextID); err != nil {
			return err
		}
		if plan.context.OriginKind == nil {
			return ErrNotLinked
		}
		var ok bool
		if plan.source, ok = s.links.Find(*plan.context.OriginKind); !ok {
			return ErrNotLinked
		}
		if plan.platform, err = s.platforms.ByID(ctx, tx, plan.context.PlatformID); err != nil {
			return err
		}
		if plan.key, err = s.keys.Current(ctx, tx); err != nil {
			return err
		}
		tenant, err := tenantSource(ctx, tx, s.values)
		if err != nil {
			return err
		}
		plan.mode = effectiveSource(tenant, plan.context)
		if plan.context.LinkedBy != nil {
			plan.actor = *plan.context.LinkedBy
		}
		return nil
	})
	return plan, err
}

func (s SyncService) run(ctx context.Context, plan syncPlan) (SyncReport, error) {
	var report SyncReport
	if plan.context.MembersURL != nil {
		if err := s.syncMembers(ctx, plan, &report); err != nil {
			return report, err
		}
	}
	if plan.context.LineItemsURL == nil {
		return report, nil
	}
	links, err := s.syncWorks(ctx, plan, &report)
	if err != nil {
		return report, err
	}
	if plan.has(ScopeScores) {
		return report, s.syncScores(ctx, plan, links, &report)
	}
	return report, nil
}

func (s SyncService) syncMembers(ctx context.Context, plan syncPlan, report *SyncReport) error {
	token, err := s.client.Token(ctx, plan.platform, plan.key, []string{ScopeMembers})
	if err != nil {
		return err
	}
	list, err := s.client.Members(ctx, token, *plan.context.MembersURL)
	if err != nil {
		return err
	}
	type entry struct {
		userID string
		member Member
	}
	var entries []entry
	var leads, learners, ids []string
	for _, m := range list {
		if !m.Active() || m.UserID == "" {
			continue
		}
		userID, err := s.accounts.Resolve(ctx, plan.tenantID, plan.platform, m.Person())
		if err != nil {
			return err
		}
		e := entry{userID: userID, member: Member{Roles: m.Roles, Instructor: m.Instructor()}}
		entries = append(entries, e)
		ids = append(ids, userID)
		if e.member.Instructor {
			leads = append(leads, userID)
		} else {
			learners = append(learners, userID)
		}
	}
	report.Members = len(entries)
	return s.tx.Run(ctx, plan.tenantID, func(tx pgx.Tx) error {
		for _, e := range entries {
			if err := s.members.Record(ctx, tx, plan.context.ID, e.userID, e.member); err != nil {
				return err
			}
		}
		if err := s.members.Keep(ctx, tx, plan.context.ID, nonNil(ids)); err != nil {
			return err
		}
		var err error
		report.Enrolled, err = plan.source.Enroll(ctx, tx, plan.actor, *plan.context.OriginID, leads, learners)
		return err
	})
}

func (s SyncService) syncWorks(ctx context.Context, plan syncPlan, report *SyncReport) ([]LineItemLink, error) {
	scopes := []string{}
	for _, scope := range []string{ScopeLineItems, ScopeLineItemsReadOnly} {
		if plan.has(scope) {
			scopes = append(scopes, scope)
		}
	}
	var links []LineItemLink
	err := s.tx.Run(ctx, plan.tenantID, func(tx pgx.Tx) error {
		var err error
		links, err = s.lineItems.ByContext(ctx, tx, plan.context.ID)
		return err
	})
	if err != nil || len(scopes) == 0 {
		return links, err
	}
	token, err := s.client.Token(ctx, plan.platform, plan.key, scopes)
	if err != nil {
		return nil, err
	}
	items, err := s.client.LineItems(ctx, token, *plan.context.LineItemsURL)
	if err != nil {
		return nil, err
	}
	switch {
	case plan.mode == SourceLMS:
		return s.importWorks(ctx, plan, items, links, report)
	case plan.has(ScopeLineItems):
		return s.exportWorks(ctx, plan, token, items, links, report)
	default:
		return links, nil
	}
}

func (s SyncService) importWorks(ctx context.Context, plan syncPlan, items []LineItem, links []LineItemLink, report *SyncReport) ([]LineItemLink, error) {
	byURL := map[string]LineItemLink{}
	for _, l := range links {
		byURL[l.URL] = l
	}
	out := make([]LineItemLink, 0, len(items))
	for _, item := range items {
		if item.ID == "" || item.ScoreMaximum <= 0 {
			continue
		}
		if l, ok := byURL[item.ID]; ok && item.Tag == lodestarTag {
			out = append(out, l)
			continue
		}
		err := s.tx.Run(ctx, plan.tenantID, func(tx pgx.Tx) error {
			ref, err := plan.source.PutWork(ctx, tx, plan.actor, *plan.context.OriginID, byURL[item.ID].WorkID, module.ExternalWork{
				Title: item.Label, Max: item.ScoreMaximum, DueAt: item.EndDateTime,
			})
			if err != nil {
				return err
			}
			l, err := s.lineItems.Put(ctx, tx, plan.context.ID, item.ID, ref.Kind, ref.ID)
			out = append(out, l)
			return err
		})
		if err != nil {
			return nil, fmt.Errorf("import %q: %w", item.Label, err)
		}
		report.Works++
	}
	return out, nil
}

func (s SyncService) exportWorks(ctx context.Context, plan syncPlan, token string, items []LineItem, links []LineItemLink, report *SyncReport) ([]LineItemLink, error) {
	works, err := s.works(ctx, plan)
	if err != nil {
		return nil, err
	}
	mapped := map[string]bool{}
	for _, l := range links {
		mapped[l.WorkKind+" "+l.WorkID] = true
	}
	byResource := map[string]LineItem{}
	for _, item := range items {
		if item.ResourceID != "" {
			byResource[item.ResourceID] = item
		}
	}
	for _, w := range works {
		if mapped[w.Kind+" "+w.ID] || w.Max <= 0 {
			continue
		}
		item, ok := byResource[w.ID]
		if !ok {
			if item, err = s.client.CreateLineItem(ctx, token, *plan.context.LineItemsURL, LineItem{
				Label: w.Title, ScoreMaximum: w.Max, ResourceID: w.ID, Tag: lodestarTag,
			}); err != nil {
				return nil, err
			}
		}
		err := s.tx.Run(ctx, plan.tenantID, func(tx pgx.Tx) error {
			l, err := s.lineItems.Put(ctx, tx, plan.context.ID, item.ID, w.Kind, w.ID)
			links = append(links, l)
			return err
		})
		if err != nil {
			return nil, err
		}
		report.Works++
	}
	return links, nil
}

func (s SyncService) syncScores(ctx context.Context, plan syncPlan, links []LineItemLink, report *SyncReport) error {
	if len(links) == 0 {
		return nil
	}
	works, err := s.works(ctx, plan)
	if err != nil {
		return err
	}
	var subjects map[string]string
	err = s.tx.Run(ctx, plan.tenantID, func(tx pgx.Tx) error {
		subjects, err = s.identities.Subjects(ctx, tx, plan.platform.ID)
		return err
	})
	if err != nil {
		return err
	}
	byWork := map[string]module.ScoredWork{}
	for _, w := range works {
		byWork[w.Kind+" "+w.ID] = w
	}
	token, err := s.client.Token(ctx, plan.platform, plan.key, []string{ScopeScores})
	if err != nil {
		return err
	}
	for _, l := range links {
		w, ok := byWork[l.WorkKind+" "+l.WorkID]
		if !ok {
			continue
		}
		for _, r := range w.Results {
			subject, ok := subjects[r.UserID]
			if !ok || r.Max <= 0 {
				continue
			}
			if err := s.pushScore(ctx, plan, token, l, subject, r, report); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s SyncService) pushScore(ctx context.Context, plan syncPlan, token string, l LineItemLink, subject string, r module.WorkResult, report *SyncReport) error {
	var changed bool
	err := s.tx.Run(ctx, plan.tenantID, func(tx pgx.Tx) error {
		var err error
		changed, err = s.sent.Changed(ctx, tx, l.ID, r.UserID, r.Score, r.Max)
		return err
	})
	if err != nil || !changed {
		return err
	}
	at := r.At
	if at.IsZero() {
		at = time.Now()
	}
	err = s.client.PostScore(ctx, token, l.URL, Score{
		UserID: subject, ScoreGiven: r.Score, ScoreMaximum: r.Max,
		ActivityProgress: "Completed", GradingProgress: "FullyGraded", Timestamp: at,
	})
	if err != nil {
		return err
	}
	report.Scores++
	return s.tx.Run(ctx, plan.tenantID, func(tx pgx.Tx) error {
		return s.sent.Sent(ctx, tx, l.ID, r.UserID, r.Score, r.Max)
	})
}

func (s SyncService) works(ctx context.Context, plan syncPlan) ([]module.ScoredWork, error) {
	source, ok := s.scores.Find(plan.source.Kind())
	if !ok {
		return nil, nil
	}
	var works []module.ScoredWork
	err := s.tx.Run(ctx, plan.tenantID, func(tx pgx.Tx) error {
		var err error
		works, err = source.ScoredWorks(ctx, tx, *plan.context.OriginID)
		return err
	})
	return works, err
}
