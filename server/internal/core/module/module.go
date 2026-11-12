package module

import "net/http"

type Module interface {
	Manifest() Manifest
	Mount(mux *http.ServeMux)
}
