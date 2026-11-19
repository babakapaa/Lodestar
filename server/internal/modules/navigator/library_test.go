package navigator_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/modules/navigator"
)

const openAlexReply = `{"results":[
	{"display_name":"Photosynthesis: light and dark reactions","doi":"https://doi.org/10.1000/xyz123","publication_year":2021,
	 "authorships":[{"author":{"display_name":"A. Smith"}},{"author":{"display_name":"B. Jones"}}],
	 "primary_location":{"source":{"display_name":"Plant Biology"},"landing_page_url":"https://example.org/a"}},
	{"display_name":"","doi":null,"publication_year":2020,"authorships":[],"primary_location":{}}
]}`

type redirect struct {
	target *url.URL
}

func (r redirect) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	req.URL.Scheme, req.URL.Host = r.target.Scheme, r.target.Host
	return http.DefaultTransport.RoundTrip(req)
}

func libraryWorld(t *testing.T, on bool) (*world, *[]url.Values) {
	var queries []url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		queries = append(queries, r.URL.Query())
		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(openAlexReply))
	}))
	t.Cleanup(srv.Close)
	target, _ := url.Parse(srv.URL)
	w := newWorld(t)
	w.kernel.Outbound = &http.Client{Transport: redirect{target: target}}
	w.kernel.Settings.(*fakeSettings).values[navigator.LibrarySetting] = on
	w.navigator = navigator.NewService(w.kernel)
	return w, &queries
}

func TestNavigatorFindsArticlesWhenCourseHasNothing(t *testing.T) {
	w, queries := libraryWorld(t, true)
	student := w.user("aigerim")
	w.enrol(module.ModeStudy, student)
	w.model.reply = "photosynthesis light reactions"
	if _, _, err := w.navigator.Create(w.ctx, student, navigator.DialogInput{ID: dialogID, ScopeKind: "assignment", ScopeID: assignmentID}); err != nil {
		t.Fatalf("create: %v", err)
	}
	a, _, err := w.navigator.Ask(w.ctx, student, dialogID, navigator.MessageInput{Key: "q1", Text: "Где почитать про фотосинтез?", Intent: navigator.IntentFind}, nil)
	if err != nil {
		t.Fatalf("ask: %v", err)
	}
	if len(*queries) != 1 || (*queries)[0].Get("search") != "photosynthesis light reactions" || (*queries)[0].Get("filter") != "has_doi:true" {
		t.Fatalf("queries = %v", *queries)
	}
	got := a.Reply.Articles
	if len(got) != 1 || got[0].DOI != "10.1000/xyz123" || got[0].URL != "https://doi.org/10.1000/xyz123" || got[0].Venue != "Plant Biology" || len(got[0].Authors) != 2 {
		t.Fatalf("articles = %+v", got)
	}
	sum, err := w.navigator.Help(w.ctx, student, "assignment", assignmentID)
	if err != nil || len(sum.Articles) != 1 || sum.Articles[0] != "A. Smith. Photosynthesis: light and dark reactions. 2021. DOI: 10.1000/xyz123" {
		t.Fatalf("help = %+v, %v", sum, err)
	}
	list, err := w.navigator.Journal(w.ctx, student, "assignment", assignmentID)
	if err != nil || len(list) != 1 || len(list[0].Sources) != 1 || list[0].Sources[0].URL != "https://doi.org/10.1000/xyz123" {
		t.Fatalf("journal = %+v, %v", list, err)
	}
	if _, _, err := w.navigator.Ask(w.ctx, student, dialogID, navigator.MessageInput{Key: "q2", Text: "Объясни фотосинтез", Intent: navigator.IntentUnderstand}, nil); err != nil {
		t.Fatalf("ask: %v", err)
	}
	if len(*queries) != 1 {
		t.Fatalf("searched on a question step: %v", *queries)
	}
}

func TestUniversityCanTurnOpenLibraryOff(t *testing.T) {
	w, queries := libraryWorld(t, false)
	student := w.user("aigerim")
	w.enrol(module.ModeStudy, student)
	if _, _, err := w.navigator.Create(w.ctx, student, navigator.DialogInput{ID: dialogID, ScopeKind: "assignment", ScopeID: assignmentID}); err != nil {
		t.Fatalf("create: %v", err)
	}
	a, _, err := w.navigator.Ask(w.ctx, student, dialogID, navigator.MessageInput{Key: "q1", Text: "Где почитать про фотосинтез?", Intent: navigator.IntentFind}, nil)
	if err != nil || len(a.Reply.Articles) != 0 || len(*queries) != 0 {
		t.Fatalf("articles = %+v, queries = %v, %v", a.Reply.Articles, *queries, err)
	}
}
