package lti

import "strings"

const (
	MessageResourceLink = "LtiResourceLinkRequest"
	MessageDeepLinking  = "LtiDeepLinkingRequest"
	messageDeepAnswer   = "LtiDeepLinkingResponse"
	customKind          = "lodestar_kind"
	customID            = "lodestar_id"
	Version             = "1.3.0"
	roleInstructor      = "membership#Instructor"
	roleAdministrator   = "membership#Administrator"
	roleContentDev      = "membership#ContentDeveloper"
)

type LaunchClaims struct {
	AuthorizedBy string            `json:"azp"`
	Email        string            `json:"email"`
	Name         string            `json:"name"`
	GivenName    string            `json:"given_name"`
	FamilyName   string            `json:"family_name"`
	MessageType  string            `json:"https://purl.imsglobal.org/spec/lti/claim/message_type"`
	Version      string            `json:"https://purl.imsglobal.org/spec/lti/claim/version"`
	DeploymentID string            `json:"https://purl.imsglobal.org/spec/lti/claim/deployment_id"`
	Roles        []string          `json:"https://purl.imsglobal.org/spec/lti/claim/roles"`
	Context      ContextClaim      `json:"https://purl.imsglobal.org/spec/lti/claim/context"`
	Names        NamesClaim        `json:"https://purl.imsglobal.org/spec/lti-nrps/claim/namesroleservice"`
	Grades       EndpointClaim     `json:"https://purl.imsglobal.org/spec/lti-ags/claim/endpoint"`
	DeepLinking  *DeepLinkingClaim `json:"https://purl.imsglobal.org/spec/lti-dl/claim/deep_linking_settings"`
	Custom       map[string]string `json:"https://purl.imsglobal.org/spec/lti/claim/custom"`
}

func (c LaunchClaims) DisplayName() string {
	if name := strings.TrimSpace(c.Name); name != "" {
		return name
	}
	return strings.TrimSpace(c.GivenName + " " + c.FamilyName)
}

func (c LaunchClaims) Instructor() bool {
	for _, role := range c.Roles {
		if strings.HasSuffix(role, roleInstructor) || strings.HasSuffix(role, roleAdministrator) || strings.HasSuffix(role, roleContentDev) {
			return true
		}
	}
	return false
}

func (c LaunchClaims) Valid() bool {
	if c.Version != Version || c.DeploymentID == "" || c.Context.ID == "" {
		return false
	}
	switch c.MessageType {
	case MessageResourceLink:
		return true
	case MessageDeepLinking:
		return c.DeepLinking != nil && c.DeepLinking.ReturnURL != "" && c.DeepLinking.acceptsLinks()
	}
	return false
}

func (c LaunchClaims) Target() (string, string, bool) {
	kind, id := c.Custom[customKind], c.Custom[customID]
	return kind, id, kind != "" && id != ""
}
