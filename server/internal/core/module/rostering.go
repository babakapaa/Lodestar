package module

type Rostering interface {
	RosterSources(k Kernel) []RosterSource
}
