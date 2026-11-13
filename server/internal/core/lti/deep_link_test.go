package lti_test

import (
	"io"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"

	"github.com/babakapaa/lodestar/server/internal/core/lti"
)

var jwtField = regexp.MustCompile(`name="JWT" value="([^"]+)"`)

func (h *harness) deepLinkClaims(nonce, subject, email string, roles ...string) map[string]any {
	c := h.claims(nonce, subject, email, roles...)
	c["https://purl.imsglobal.org/spec/lti/claim/message_type"] = lti.MessageDeepLinking
	c["https://purl.imsglobal.org/spec/lti-dl/claim/deep_linking_settings"] = map[string]any{
		"deep_link_return_url": h.lms.url + "/deep-links",
		"accept_types":         []string{"ltiResourceLink"},
		"accept_multiple":      false,
		"data":                 "opaque",
	}
	return c
}

func TestDeepLinkingReturnsChosenCourse(t *testing.T) {
	h := newHarness(t)

	q := h.login()
	result := h.launch(q.Get("state"), h.deepLinkClaims(q.Get("nonce"), "t-1", "teacher@kaznu.kz", instructorRole))
	teacher, next := h.claim(result.Get("ticket")), result.Get("next")
	h.links.people(h.me(teacher).ID, "")
	if !strings.HasPrefix(next, "/lti/deep-links/") {
		t.Fatalf("next = %s", next)
	}
	id := strings.TrimPrefix(next, "/lti/deep-links/")

	var d lti.DeepLink
	if status := h.call(http.MethodGet, "/v1/lti/deep-links/"+id, teacher, nil, &d); status != http.StatusOK {
		t.Fatalf("deep link status %d", status)
	}
	if d.AcceptMultiple || d.Platform != "Moodle" || len(d.Targets) != 1 || d.Targets[0].ID != "bio" {
		t.Fatalf("deep link = %+v", d)
	}
	two := lti.DeepLinkInput{Items: []lti.DeepLinkItem{{Kind: "study.course", ID: "bio"}, {Kind: "study.course", ID: "bio"}}}
	if status := h.call(http.MethodPost, "/v1/lti/deep-links/"+id, teacher, two, nil); status != http.StatusBadRequest {
		t.Fatalf("two items status %d", status)
	}
	student, _ := h.enter("s-1", "student@kaznu.kz", learnerRole)
	one := lti.DeepLinkInput{Items: []lti.DeepLinkItem{{Kind: "study.course", ID: "bio"}}}
	if status := h.call(http.MethodPost, "/v1/lti/deep-links/"+id, student, one, nil); status != http.StatusNotFound {
		t.Fatalf("student answer status %d", status)
	}
	var answer lti.DeepLinkAnswer
	if status := h.call(http.MethodPost, "/v1/lti/deep-links/"+id, teacher, one, &answer); status != http.StatusOK {
		t.Fatalf("answer status %d", status)
	}

	resp, err := h.client.Get(answer.FormURL)
	if err != nil {
		t.Fatalf("form: %v", err)
	}
	page, _ := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK || !strings.Contains(string(page), `action="`+h.lms.url+`/deep-links"`) {
		t.Fatalf("form %d: %s", resp.StatusCode, page)
	}
	match := jwtField.FindStringSubmatch(string(page))
	if match == nil {
		t.Fatalf("no JWT in %s", page)
	}
	token, err := jwt.ParseSigned(match[1], []jose.SignatureAlgorithm{jose.RS256})
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	var claims struct {
		Issuer   string `json:"iss"`
		Audience string `json:"aud"`
		Type     string `json:"https://purl.imsglobal.org/spec/lti/claim/message_type"`
		Data     string `json:"https://purl.imsglobal.org/spec/lti-dl/claim/data"`
		Items    []struct {
			Type   string            `json:"type"`
			Title  string            `json:"title"`
			Custom map[string]string `json:"custom"`
		} `json:"https://purl.imsglobal.org/spec/lti-dl/claim/content_items"`
	}
	if err := token.UnsafeClaimsWithoutVerification(&claims); err != nil {
		t.Fatalf("claims: %v", err)
	}
	if claims.Issuer != clientID || claims.Audience != lmsIss || claims.Type != "LtiDeepLinkingResponse" || claims.Data != "opaque" {
		t.Fatalf("response = %+v", claims)
	}
	if len(claims.Items) != 1 || claims.Items[0].Type != "ltiResourceLink" || claims.Items[0].Custom["lodestar_id"] != "bio" {
		t.Fatalf("items = %+v", claims.Items)
	}
	again, err := h.client.Get(answer.FormURL)
	if err != nil {
		t.Fatalf("form again: %v", err)
	}
	_ = again.Body.Close()
	if again.StatusCode != http.StatusNotFound {
		t.Fatalf("form reused = %d", again.StatusCode)
	}

	q = h.login()
	c := h.claims(q.Get("nonce"), "s-1", "student@kaznu.kz", learnerRole)
	c["https://purl.imsglobal.org/spec/lti/claim/context"] = map[string]any{"id": "c-202", "title": "Другой курс"}
	c["https://purl.imsglobal.org/spec/lti/claim/custom"] = map[string]string{"lodestar_kind": "study.course", "lodestar_id": "bio"}
	if next := h.launch(q.Get("state"), c).Get("next"); next != "/study/courses/bio" {
		t.Fatalf("custom launch next = %s", next)
	}
}

func TestLearnerCannotDeepLink(t *testing.T) {
	h := newHarness(t)
	q := h.login()
	result := h.launch(q.Get("state"), h.deepLinkClaims(q.Get("nonce"), "s-1", "student@kaznu.kz", learnerRole))
	if result.Get("error") == "" || result.Get("ticket") != "" {
		t.Fatalf("learner deep link = %v", result)
	}
}
