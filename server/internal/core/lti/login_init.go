package lti

type LoginInit struct {
	TenantSlug   string
	Issuer       string
	ClientID     string
	DeploymentID string
	LoginHint    string
	MessageHint  string
}
