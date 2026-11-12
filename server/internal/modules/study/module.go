package study

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/files"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/core/privacy"
	"github.com/babakapaa/lodestar/server/internal/core/search"
)

const Name = "study"

type Module struct{}

func (Module) Manifest() module.Manifest {
	return module.Manifest{Name: Name, Version: "0.1.0", Optional: true}
}

func (Module) Mount(mux *http.ServeMux, k module.Kernel) {
	Routes{
		Service:     NewService(k.Tenants, k.Access),
		Assignments: NewAssignmentService(k.Tenants, k.Access, k.Notifier),
	}.Mount(mux)
}

func (Module) FileLinks(k module.Kernel) []files.Link {
	return []files.Link{FileLink{gate: CourseGate{access: k.Access}}}
}

func (Module) Permissions() []string {
	return []string{string(ManageCourses)}
}

func (Module) SearchSources() []search.Source {
	return []search.Source{CoursesSource{}}
}

func (Module) PrivacySections() []privacy.Section {
	return []privacy.Section{PrivacySection{}, WorkPrivacySection{}}
}
