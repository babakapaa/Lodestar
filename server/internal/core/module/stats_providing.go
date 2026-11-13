package module

type StatsProviding interface {
	StatsSources(k Kernel) []StatsSource
}
