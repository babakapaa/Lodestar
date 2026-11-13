package module

type GPUGranting interface {
	GPUSources(k Kernel) []GPUSource
}
