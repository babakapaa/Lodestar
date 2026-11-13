package league

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/core/privacy"
)

const Name = "league"

type Module struct{}

func (Module) Manifest() module.Manifest {
	return module.Manifest{Name: Name, Version: "0.1.0", Optional: true}
}

func (Module) Mount(mux *http.ServeMux, k module.Kernel) {
	Routes{Leagues: NewService(k)}.Mount(mux)
}

func (Module) PrivacySections(module.Kernel) []privacy.Section {
	return []privacy.Section{PrivacySection{}}
}

func (Module) Workers(k module.Kernel) []module.Worker {
	return []module.Worker{NewLeagueWorker(k)}
}
