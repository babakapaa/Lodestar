package auth

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/secretbox"
)

type Settings struct {
	Box        secretbox.Box
	Lifetimes  Lifetimes
	Redirects  []string
	PublicURL  string
	HTTPClient *http.Client
}
