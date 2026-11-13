package attendance

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/core/privacy"
)

const Name = "attendance"

type Module struct{}

func (Module) Manifest() module.Manifest {
	return module.Manifest{Name: Name, Version: "0.1.0", Optional: true}
}

func (Module) Mount(mux *http.ServeMux, k module.Kernel) {
	Routes{Service: NewService(k)}.Mount(mux)
}

func (Module) PrivacySections(module.Kernel) []privacy.Section {
	return []privacy.Section{PrivacySection{}}
}

func (Module) StatsSources(module.Kernel) []module.StatsSource {
	return []module.StatsSource{AttendanceStats{}}
}
