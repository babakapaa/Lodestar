package auth

import (
	"net/url"
	"strings"
)

func CallbackURL(publicURL, tenantSlug, providerKey string) string {
	return strings.TrimRight(publicURL, "/") + "/v1/auth/" + url.PathEscape(tenantSlug) + "/" + url.PathEscape(providerKey) + "/callback"
}
