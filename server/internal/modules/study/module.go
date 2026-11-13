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
		Assignments: NewAssignmentService(k.Tenants, k.Access, k.Notifier).WithTeams(k.Teams),
		Schedule:    NewScheduleService(k),
		Defenses:    NewDefenseService(k),
		Gradebook:   NewGradebookService(k.Tenants, k.Access),
		Imports:     NewStudentImportService(k.Tenants, k.Access),
	}.Mount(mux)
}

func (Module) Workers(k module.Kernel) []module.Worker {
	return []module.Worker{NewDeadlineReminder(k), NewDefenseReminder(k)}
}

func (Module) FileLinks(k module.Kernel) []files.Link {
	return []files.Link{FileLink{gate: CourseGate{access: k.Access}}}
}

func (Module) Permissions() []string {
	return []string{string(ManageCourses), string(ManageSchedule)}
}

func (Module) SearchSources() []search.Source {
	return []search.Source{CoursesSource{}}
}

func (Module) PrivacySections(module.Kernel) []privacy.Section {
	return []privacy.Section{PrivacySection{}, WorkPrivacySection{}, SearchMapPrivacySection{}, DefensePrivacySection{}}
}

func (Module) Briefers(k module.Kernel) []module.Briefer {
	return []module.Briefer{
		AssignmentBriefer{service: NewAssignmentService(k.Tenants, k.Access, k.Notifier)},
		CourseBriefer{gate: CourseGate{access: k.Access}},
	}
}

func (Module) RosterSources(module.Kernel) []module.RosterSource {
	return []module.RosterSource{CourseRosters{}}
}

func (Module) TeamScopeSources(k module.Kernel) []module.TeamScopeSource {
	return []module.TeamScopeSource{AssignmentTeamScopes{work: NewAssignmentService(k.Tenants, k.Access, k.Notifier)}}
}

func (Module) ScoreSources(module.Kernel) []module.ScoreSource {
	return []module.ScoreSource{CourseScores{}}
}

func (Module) LinkSources(k module.Kernel) []module.LinkSource {
	return []module.LinkSource{CourseLinks{gate: CourseGate{access: k.Access}}}
}

func (Module) LessonSource(k module.Kernel) module.LessonSource {
	return CourseLessons{settings: k.Settings}
}

func (Module) StatsSources(module.Kernel) []module.StatsSource {
	return []module.StatsSource{StudyStats{}}
}

func (Module) PilotSources(k module.Kernel) []module.PilotSource {
	return []module.PilotSource{CoursePilot{settings: k.Settings}}
}

func (Module) ReadinessSources(k module.Kernel) []module.ReadinessSource {
	return []module.ReadinessSource{CourseReadiness{settings: k.Settings}}
}
