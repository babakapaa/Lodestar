package module

type HardwareProviding interface {
	HardwareSources(k Kernel) []HardwareSource
}
