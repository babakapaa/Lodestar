package module

type ReadinessProviding interface {
	ReadinessSources(k Kernel) []ReadinessSource
}
