package module

import "net/http"

type Gate interface {
	Wrap(name string, h http.Handler) http.Handler
}
