package navigator

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/core/privacy"
)

const Name = "navigator"

type Module struct{}

func (Module) Manifest() module.Manifest {
	return module.Manifest{Name: Name, Version: "0.1.0", Optional: true}
}

func (Module) Mount(mux *http.ServeMux, k module.Kernel) {
	Routes{Service: NewService(k), Materials: NewMaterialService(k), Topics: NewTopicService(k)}.Mount(mux)
}

func (Module) Workers(k module.Kernel) []module.Worker {
	if k.Embedder == nil {
		return nil
	}
	return []module.Worker{NewIndexer(k)}
}

func (Module) PrivacySections(k module.Kernel) []privacy.Section {
	return []privacy.Section{
		PrivacySection{service: NewService(k)},
		JournalPrivacySection{journal: JournalStore{sealer: Sealer{box: k.Secrets}}},
	}
}

func (Module) StatsSources(module.Kernel) []module.StatsSource {
	return []module.StatsSource{NavigatorStats{}}
}
