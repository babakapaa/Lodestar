package module

type Working interface {
	Workers(k Kernel) []Worker
}
