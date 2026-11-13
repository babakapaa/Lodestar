package lti

import "net/url"

type ToolURLs struct {
	LoginURL  string `json:"loginUrl"`
	LaunchURL string `json:"launchUrl"`
	JWKSURL   string `json:"jwksUrl"`
}

func NewToolURLs(publicURL, tenantSlug string) ToolURLs {
	base := publicURL + "/v1/lti/tenants/" + url.PathEscape(tenantSlug)
	return ToolURLs{LoginURL: base + "/login", LaunchURL: base + "/launch", JWKSURL: base + "/jwks"}
}
