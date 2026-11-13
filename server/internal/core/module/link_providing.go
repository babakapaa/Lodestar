package module

type LinkProviding interface {
	LinkSources(k Kernel) []LinkSource
}
