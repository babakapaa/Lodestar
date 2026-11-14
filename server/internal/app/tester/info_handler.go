package tester

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/app/demo"
	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type InfoHandler struct{}

func (InfoHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	info := Info{Tenant: demo.TenantSlug, People: []Person{}}
	for _, p := range demo.People() {
		info.People = append(info.People, Person{Key: p.Subject, Name: p.Name, Title: p.Title, Email: p.Email})
	}
	jsonhttp.Write(w, http.StatusOK, info)
}
