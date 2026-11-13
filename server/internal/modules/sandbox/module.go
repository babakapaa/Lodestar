package sandbox

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/module"
)

const Name = "sandbox"

type Module struct{}

func (Module) Manifest() module.Manifest {
	return module.Manifest{Name: Name, Version: "0.1.0", Optional: true}
}

func (Module) Mount(mux *http.ServeMux, k module.Kernel) {
	Routes{Quotas: NewService(k)}.Mount(mux)
}

func (Module) MountOpen(mux *http.ServeMux, k module.Kernel) {
	Routes{Quotas: NewService(k)}.MountOpen(mux)
}

func (Module) Workers(k module.Kernel) []module.Worker {
	return []module.Worker{NewRunWorker(k)}
}

func (Module) StatsSources(module.Kernel) []module.StatsSource {
	return []module.StatsSource{SandboxStats{}}
}
