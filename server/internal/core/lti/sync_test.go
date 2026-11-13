package lti_test

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/lti"
)

func (h *harness) linkedCourse() (string, string) {
	h.t.Helper()
	teacher, next := h.enter("t-1", "teacher@kaznu.kz", instructorRole)
	h.links.people(h.me(teacher).ID, "")
	contextID := strings.TrimPrefix(next, "/lti/contexts/")
	if status := h.call(http.MethodPut, "/v1/lti/contexts/"+contextID+"/link", teacher, lti.LinkInput{Kind: "study.course", ID: "bio"}, nil); status != http.StatusOK {
		h.t.Fatalf("link status %d", status)
	}
	h.lms.mu.Lock()
	h.lms.members = []lti.NamesMember{
		{UserID: "t-1", Status: "Active", Roles: []string{instructorRole}, Email: "teacher@kaznu.kz", Name: "Преподаватель"},
		{UserID: "s-1", Status: "Active", Roles: []string{learnerRole}, Email: "s1@kaznu.kz", GivenName: "Айгерим", FamilyName: "Омарова"},
		{UserID: "s-2", Status: "Active", Roles: []string{learnerRole}, Email: "s2@kaznu.kz", Name: "Данияр"},
		{UserID: "s-3", Status: "Inactive", Roles: []string{learnerRole}, Email: "s3@kaznu.kz"},
	}
	h.lms.mu.Unlock()
	return teacher, contextID
}

func (h *harness) sync(bearer, contextID string) lti.SyncReport {
	h.t.Helper()
	var report lti.SyncReport
	if status := h.call(http.MethodPost, "/v1/lti/contexts/"+contextID+"/sync", bearer, nil, &report); status != http.StatusOK {
		h.t.Fatalf("sync status %d", status)
	}
	return report
}

func TestSyncSendsWorksAndScoresToLMS(t *testing.T) {
	h := newHarness(t)
	teacher, contextID := h.linkedCourse()

	report := h.sync(teacher, contextID)
	if report != (lti.SyncReport{Members: 3, Enrolled: 3, Works: 1, Scores: 2}) {
		t.Fatalf("report = %+v", report)
	}
	h.links.mu.Lock()
	leads, learners := len(h.links.leads), len(h.links.learners)
	h.links.mu.Unlock()
	if leads != 1 || learners != 2 {
		t.Fatalf("enrolled leads %d learners %d", leads, learners)
	}
	h.lms.mu.Lock()
	items, scores := h.lms.items, h.lms.scores
	h.lms.mu.Unlock()
	if len(items) != 1 || items[0].ResourceID != "essay" || items[0].ScoreMaximum != 10 || items[0].Tag != "lodestar" {
		t.Fatalf("items = %+v", items)
	}
	if len(scores) != 2 || scores[0].ScoreGiven != 8 || scores[0].GradingProgress != "FullyGraded" || !strings.HasPrefix(scores[0].UserID, "s-") {
		t.Fatalf("scores = %+v", scores)
	}

	again := h.sync(teacher, contextID)
	if again.Works != 0 || again.Scores != 0 || again.Members != 3 {
		t.Fatalf("second report = %+v", again)
	}
	var view lti.ContextView
	h.call(http.MethodGet, "/v1/lti/contexts/"+contextID, teacher, nil, &view)
	if view.SyncedAt == nil || view.SyncError != nil || view.Source != lti.SourceLodestar || !view.CanSync {
		t.Fatalf("view = %+v", view)
	}

	student, _ := h.enter("s-1", "s1@kaznu.kz", learnerRole)
	if status := h.call(http.MethodPost, "/v1/lti/contexts/"+contextID+"/sync", student, nil, nil); status != http.StatusForbidden {
		t.Fatalf("student sync status %d", status)
	}
}

func TestSyncTakesWorksFromLMS(t *testing.T) {
	h := newHarness(t)
	h.values.set(lti.SourceLMS)
	teacher, contextID := h.linkedCourse()
	h.lms.mu.Lock()
	h.lms.items = []lti.LineItem{{ID: h.lms.url + "/lineitems/7", Label: "Эссе", ScoreMaximum: 20}}
	h.lms.mu.Unlock()

	if report := h.sync(teacher, contextID); report.Works != 1 {
		t.Fatalf("report = %+v", report)
	}
	h.sync(teacher, contextID)
	h.links.mu.Lock()
	works, updated := h.links.works, h.links.updated
	h.links.mu.Unlock()
	if len(works) != 1 || works[0].Title != "Эссе" || works[0].Max != 20 || len(updated) != 1 || updated[0] != "w-Эссе" {
		t.Fatalf("works %+v updated %+v", works, updated)
	}
	h.lms.mu.Lock()
	created := len(h.lms.items)
	h.lms.mu.Unlock()
	if created != 1 {
		t.Fatalf("lms items = %d", created)
	}
}

func TestSyncRecordsFailure(t *testing.T) {
	h := newHarness(t)
	teacher, contextID := h.linkedCourse()
	_, err := h.routes.Service.UpdatePlatform(h.ctx, h.admin, h.platform.ID, lti.PlatformInput{
		Name: "Moodle", Issuer: lmsIss, ClientID: clientID, AuthURL: lmsIss + "/auth",
		TokenURL: h.lms.url + "/missing", JWKSURL: h.lms.url + "/jwks", DeploymentIDs: []string{deployID},
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	status := h.call(http.MethodPost, "/v1/lti/contexts/"+contextID+"/sync", teacher, nil, nil)
	if status != http.StatusBadGateway {
		t.Fatalf("sync status %d", status)
	}
	var view lti.ContextView
	h.call(http.MethodGet, "/v1/lti/contexts/"+contextID, teacher, nil, &view)
	if view.SyncError == nil || !strings.Contains(*view.SyncError, "answered 401") {
		t.Fatalf("view = %+v", view)
	}
}

func TestAssignmentSourceChoice(t *testing.T) {
	h := newHarness(t)
	teacher, contextID := h.linkedCourse()
	head := h.user("head", "head@kaznu.kz")
	h.links.people(h.me(teacher).ID, head.UserID)
	lms := lti.SourceLMS
	if _, err := h.routes.Service.SetSource(h.ctx, head, contextID, lti.SourceInput{Source: &lms}); !errors.Is(err, lti.ErrSourceFixed) {
		t.Fatalf("fixed err = %v", err)
	}
	h.values.set(lti.SourceUnits)
	teacherPrincipal := auth.Principal{TenantID: h.tenant, UserID: h.me(teacher).ID}
	if _, err := h.routes.Service.SetSource(h.ctx, teacherPrincipal, contextID, lti.SourceInput{Source: &lms}); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("teacher err = %v", err)
	}
	view, err := h.routes.Service.SetSource(h.ctx, head, contextID, lti.SourceInput{Source: &lms})
	if err != nil || view.Source != lti.SourceLMS || !view.CanChoose || !view.CanSync {
		t.Fatalf("view = %+v, %v", view, err)
	}
	list, err := h.routes.Service.ByOrigin(h.ctx, head, "study.course", "bio")
	if err != nil || len(list) != 1 || list[0].ID != contextID {
		t.Fatalf("by origin = %+v, %v", list, err)
	}
	stranger := h.user("stranger", "x@kaznu.kz")
	if _, err := h.routes.Service.ByOrigin(h.ctx, stranger, "study.course", "bio"); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("stranger err = %v", err)
	}
}
