package lti

import (
	"html/template"
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

var deepLinkPage = template.Must(template.New("deep-link").Parse(`<!doctype html>
<html><head><meta charset="utf-8"><title>Lodestar</title></head>
<body onload="document.forms[0].submit()">
<form method="post" action="{{.URL}}"><input type="hidden" name="JWT" value="{{.JWT}}"><noscript><button type="submit">Continue</button></noscript></form>
</body></html>`))

type DeepLinkHandler struct {
	service Service
}

func (h DeepLinkHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d, err := h.service.DeepLink(r.Context(), principal(r), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, d)
}

type AnswerDeepLinkHandler struct {
	service Service
}

func (h AnswerDeepLinkHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var in DeepLinkInput
	if err := jsonhttp.Read(w, r, &in); err != nil {
		badRequest(w)
		return
	}
	a, err := h.service.AnswerDeepLink(r.Context(), principal(r), r.PathValue("id"), in)
	if err != nil {
		writeError(w, err)
		return
	}
	jsonhttp.Write(w, http.StatusOK, a)
}

type DeepLinkReturnHandler struct {
	launch LaunchService
}

func (h DeepLinkReturnHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	returnURL, response, err := h.launch.DeepLinkForm(r.Context(), r.PathValue("tenant"), r.URL.Query().Get("code"))
	if err != nil {
		writeError(w, err)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	_ = deepLinkPage.Execute(w, struct{ URL, JWT string }{returnURL, response})
}
