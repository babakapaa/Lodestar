package module

type PilotProviding interface {
	PilotSources(k Kernel) []PilotSource
}
