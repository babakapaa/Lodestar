package module

import "net/http"

type OpenMounting interface {
	MountOpen(mux *http.ServeMux, k Kernel)
}
