package study_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/modules/study"
)

func (w *world) serve(p auth.Principal, method, path, body string) *httptest.ResponseRecorder {
	w.t.Helper()
	mux := http.NewServeMux()
	study.Routes{Service: w.study}.Mount(mux)
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req = req.WithContext(auth.WithPrincipal(req.Context(), p))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}

func TestCourseOverHTTP(t *testing.T) {
	w := newWorld(t)
	teacher := w.user("dana", "dana@kaznu.kz")

	rec := w.serve(w.admin, http.MethodPost, "/v1/study/courses", `{"title":"Биология","code":"BIO-1"}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create: %d %s", rec.Code, rec.Body)
	}
	var c study.Course
	if err := json.Unmarshal(rec.Body.Bytes(), &c); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if rec := w.serve(w.admin, http.MethodPut, "/v1/study/courses/"+c.ID+"/teachers/"+teacher.UserID, ""); rec.Code != http.StatusNoContent {
		t.Fatalf("add teacher: %d %s", rec.Code, rec.Body)
	}
	rec = w.serve(teacher, http.MethodGet, "/v1/study/courses", "")
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"role":"teacher"`) {
		t.Fatalf("mine: %d %s", rec.Code, rec.Body)
	}
	if rec := w.serve(teacher, http.MethodPost, "/v1/study/courses", `{"title":"x"}`); rec.Code != http.StatusForbidden {
		t.Fatalf("teacher create: %d %s", rec.Code, rec.Body)
	}
	if rec := w.serve(w.admin, http.MethodPost, "/v1/study/courses", `{"title":"y","code":"BIO-1"}`); rec.Code != http.StatusConflict || !strings.Contains(rec.Body.String(), study.CodeCodeTaken) {
		t.Fatalf("taken code: %d %s", rec.Code, rec.Body)
	}
	if rec := w.serve(teacher, http.MethodGet, "/v1/study/courses/"+c.ID+"/students", ""); rec.Code != http.StatusOK || rec.Body.String() != "[]\n" {
		t.Fatalf("students: %d %q", rec.Code, rec.Body)
	}
}
