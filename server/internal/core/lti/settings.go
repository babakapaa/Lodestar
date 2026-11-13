package lti

import (
	"net/http"
	"time"

	"github.com/babakapaa/lodestar/server/internal/core/access"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/platform/secretbox"
)

type Settings struct {
	Box           secretbox.Box
	Access        access.Service
	Tickets       auth.TicketStore
	Authenticator auth.Authenticator
	Links         module.LinkSources
	Scores        module.ScoreSources
	Values        Values
	HTTPClient    *http.Client
	AttemptTTL    time.Duration
	SyncEvery     time.Duration
	PublicURL     string
	AppURL        string
}
