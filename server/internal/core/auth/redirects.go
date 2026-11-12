package auth

import "slices"

type Redirects struct {
	allowed []string
}

func NewRedirects(allowed []string) Redirects {
	return Redirects{allowed: allowed}
}

func (r Redirects) Allows(uri string) bool {
	return uri != "" && slices.Contains(r.allowed, uri)
}
