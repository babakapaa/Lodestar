package livekit

import (
	"strings"

	"github.com/babakapaa/lodestar/server/internal/platform/objectstore"
)

type Config struct {
	URL       string
	APIURL    string
	APIKey    string
	APISecret string
	Storage   objectstore.Config
}

func (c Config) Enabled() bool {
	return c.URL != "" && c.APIKey != "" && c.APISecret != ""
}

func (c Config) api() string {
	if c.APIURL != "" {
		return strings.TrimRight(c.APIURL, "/")
	}
	u := strings.TrimRight(c.URL, "/")
	if rest, ok := strings.CutPrefix(u, "ws"); ok {
		return "http" + rest
	}
	return u
}
