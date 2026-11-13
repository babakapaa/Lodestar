package lti_test

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/lti"
)

const (
	instructorRole = "http://purl.imsglobal.org/vocab/lis/v2/membership#Instructor"
	learnerRole    = "http://purl.imsglobal.org/vocab/lis/v2/membership#Learner"
)

func TestLoginRedirectsToPlatform(t *testing.T) {
	h := newHarness(t)
	q := h.login()
	if q.Get("client_id") != clientID || q.Get("response_mode") != "form_post" || q.Get("lti_message_hint") != "msg" {
		t.Fatalf("auth query = %v", q)
	}
	if q.Get("redirect_uri") != h.server.URL+"/v1/lti/tenants/kaznu/launch" || q.Get("state") == "" || q.Get("nonce") == "" {
		t.Fatalf("auth query = %v", q)
	}
}

func TestLaunchSignsInAndLinksContext(t *testing.T) {
	var owner string
	h := newHarness(t, fakeLinks{owner: &owner})

	teacher, next := h.enter("t-1", "teacher@kaznu.kz", instructorRole)
	owner = h.me(teacher).ID
	if !strings.HasPrefix(next, "/lti/contexts/") {
		t.Fatalf("next = %s", next)
	}
	contextID := strings.TrimPrefix(next, "/lti/contexts/")
	var view lti.ContextView
	if status := h.call(http.MethodGet, "/v1/lti/contexts/"+contextID, teacher, nil, &view); status != http.StatusOK {
		t.Fatalf("context status %d", status)
	}
	if !view.Instructor || view.Title != "Основы биологии" || view.Platform != "Moodle" || len(view.Options) != 1 || view.Link != nil {
		t.Fatalf("view = %+v", view)
	}
	if status := h.call(http.MethodPut, "/v1/lti/contexts/"+contextID+"/link", teacher, lti.LinkInput{Kind: "study.course", ID: "bio"}, &view); status != http.StatusOK {
		t.Fatalf("link status %d", status)
	}
	if view.Link == nil || view.Link.Path != "/study/courses/bio" {
		t.Fatalf("linked view = %+v", view)
	}

	again, next := h.enter("t-1", "teacher@kaznu.kz", instructorRole)
	if next != "/study/courses/bio" || h.me(again).ID != owner {
		t.Fatalf("second launch next = %s", next)
	}

	student, next := h.enter("s-1", "student@kaznu.kz", learnerRole)
	if next != "/study/courses/bio" {
		t.Fatalf("student next = %s", next)
	}
	if status := h.call(http.MethodGet, "/v1/lti/contexts/"+contextID, student, nil, &view); status != http.StatusOK || view.Instructor || len(view.Options) != 0 {
		t.Fatalf("student view %d %+v", status, view)
	}
	if status := h.call(http.MethodDelete, "/v1/lti/contexts/"+contextID+"/link", student, nil, nil); status != http.StatusForbidden {
		t.Fatalf("student unlink status %d", status)
	}
	if status := h.call(http.MethodDelete, "/v1/lti/contexts/"+contextID+"/link", teacher, nil, &view); status != http.StatusOK || view.Link != nil {
		t.Fatalf("unlink %d %+v", status, view)
	}
}

func TestLaunchMatchesExistingUserByEmail(t *testing.T) {
	h := newHarness(t)
	existing := h.user("campus-1", "aigerim@kaznu.kz")
	bearer, _ := h.enter("lms-77", "Aigerim@kaznu.kz", learnerRole)
	if got := h.me(bearer).ID; got != existing.UserID {
		t.Fatalf("user = %s, want %s", got, existing.UserID)
	}
}

func TestLaunchRejectsBadTokens(t *testing.T) {
	h := newHarness(t)
	cases := map[string]func(map[string]any){
		"nonce":      func(c map[string]any) { c["nonce"] = "other" },
		"audience":   func(c map[string]any) { c["aud"] = "other" },
		"deployment": func(c map[string]any) { c["https://purl.imsglobal.org/spec/lti/claim/deployment_id"] = "d9" },
		"message": func(c map[string]any) {
			c["https://purl.imsglobal.org/spec/lti/claim/message_type"] = "LtiDeepLinkingRequest"
		},
		"azp": func(c map[string]any) { c["azp"] = "other" },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			q := h.login()
			claims := h.claims(q.Get("nonce"), "x-1", "x@kaznu.kz", learnerRole)
			mutate(claims)
			result := h.launch(q.Get("state"), claims)
			if result.Get("error") != "launch_failed" || result.Get("ticket") != "" {
				t.Fatalf("result = %v", result)
			}
		})
	}
	q := h.login()
	claims := h.claims(q.Get("nonce"), "x-1", "x@kaznu.kz", learnerRole)
	if result := h.launch(q.Get("state"), claims); result.Get("ticket") == "" {
		t.Fatalf("first launch = %v", result)
	}
	if result := h.launch(q.Get("state"), claims); result.Get("error") != "launch_failed" {
		t.Fatalf("replayed launch = %v", result)
	}
}

func TestTicketIsSingleUse(t *testing.T) {
	h := newHarness(t)
	q := h.login()
	result := h.launch(q.Get("state"), h.claims(q.Get("nonce"), "s-2", "s2@kaznu.kz", learnerRole))
	h.claim(result.Get("ticket"))
	status := h.call(http.MethodPost, "/v1/auth/token", "", map[string]string{"grantType": "launch_ticket", "ticket": result.Get("ticket")}, nil)
	if status != http.StatusUnauthorized && status != http.StatusBadRequest {
		t.Fatalf("second claim status %d", status)
	}
}

func TestPlatformsNeedPermission(t *testing.T) {
	h := newHarness(t)
	stranger := h.user("stranger", "stranger@kaznu.kz")
	if _, err := h.routes.Service.Platforms(h.ctx, stranger); !errors.Is(err, access.ErrForbidden) {
		t.Fatalf("stranger platforms err = %v", err)
	}
	list, err := h.routes.Service.Platforms(h.ctx, h.admin)
	if err != nil || len(list) != 1 || list[0].Issuer != lmsIss {
		t.Fatalf("platforms = %+v, %v", list, err)
	}
	_, err = h.routes.Service.CreatePlatform(h.ctx, h.admin, lti.PlatformInput{
		Name: "Copy", Issuer: lmsIss, ClientID: clientID, AuthURL: lmsIss + "/a", TokenURL: lmsIss + "/t", JWKSURL: lmsIss + "/k", DeploymentIDs: []string{"x"},
	})
	if !errors.Is(err, lti.ErrPlatformTaken) {
		t.Fatalf("duplicate err = %v", err)
	}
	_, err = h.routes.Service.CreatePlatform(h.ctx, h.admin, lti.PlatformInput{Name: "Bad", Issuer: "lms", ClientID: "c", AuthURL: "x", TokenURL: "x", JWKSURL: "x", DeploymentIDs: []string{"x"}})
	if !errors.Is(err, lti.ErrInvalidPlatform) {
		t.Fatalf("invalid err = %v", err)
	}
	tool, err := h.routes.Service.Tool(h.ctx, h.admin)
	if err != nil || tool.JWKSURL != h.server.URL+"/v1/lti/tenants/kaznu/jwks" {
		t.Fatalf("tool = %+v, %v", tool, err)
	}
	var set struct {
		Keys []map[string]any `json:"keys"`
	}
	if status := h.call(http.MethodGet, "/v1/lti/tenants/kaznu/jwks", "", nil, &set); status != http.StatusOK || len(set.Keys) != 1 || set.Keys[0]["kty"] != "RSA" || set.Keys[0]["d"] != nil {
		t.Fatalf("jwks %d %+v", status, set)
	}
	if err := h.routes.Service.DeletePlatform(h.ctx, h.admin, h.platform.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
}
