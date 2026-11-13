package lti

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
)

type Routes struct {
	Service       Service
	Launch        LaunchService
	JWKS          JWKSHandler
	Authenticator auth.Authenticator
}

func (r Routes) Mount(mux *http.ServeMux) {
	require := r.Authenticator.Require
	s := r.Service
	mux.Handle("GET /v1/lti/tenants/{tenant}/jwks", r.JWKS)
	mux.Handle("GET /v1/lti/tenants/{tenant}/login", LoginHandler{launch: r.Launch})
	mux.Handle("POST /v1/lti/tenants/{tenant}/login", LoginHandler{launch: r.Launch})
	mux.Handle("POST /v1/lti/tenants/{tenant}/launch", LaunchHandler{launch: r.Launch})
	mux.Handle("GET /v1/lti/tenants/{tenant}/deep-links/return", DeepLinkReturnHandler{launch: r.Launch})
	mux.Handle("GET /v1/lti/tool", require(ToolHandler{service: s}))
	mux.Handle("GET /v1/lti/platforms", require(PlatformsHandler{service: s}))
	mux.Handle("POST /v1/lti/platforms", require(CreatePlatformHandler{service: s}))
	mux.Handle("PUT /v1/lti/platforms/{id}", require(UpdatePlatformHandler{service: s}))
	mux.Handle("DELETE /v1/lti/platforms/{id}", require(DeletePlatformHandler{service: s}))
	mux.Handle("GET /v1/lti/contexts/{id}", require(ContextHandler{service: s}))
	mux.Handle("PUT /v1/lti/contexts/{id}/link", require(LinkHandler{service: s}))
	mux.Handle("DELETE /v1/lti/contexts/{id}/link", require(UnlinkHandler{service: s}))
	mux.Handle("PUT /v1/lti/contexts/{id}/source", require(SourceHandler{service: s}))
	mux.Handle("POST /v1/lti/contexts/{id}/sync", require(SyncHandler{service: s}))
	mux.Handle("GET /v1/lti/links/{kind}/{id}", require(OriginHandler{service: s}))
	mux.Handle("GET /v1/lti/deep-links/{id}", require(DeepLinkHandler{service: s}))
	mux.Handle("POST /v1/lti/deep-links/{id}", require(AnswerDeepLinkHandler{service: s}))
}
