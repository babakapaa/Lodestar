package app

import "net/http"

type Routes interface {
	Mount(mux *http.ServeMux)
}
