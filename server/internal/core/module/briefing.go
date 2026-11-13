package module

type Briefing interface {
	Briefers(k Kernel) []Briefer
}
