package module

type HelpProviding interface {
	HelpSource(k Kernel) HelpSource
}
