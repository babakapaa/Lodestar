package study

import "net/http"

type Routes struct {
	Service Service
}

func (r Routes) Mount(mux *http.ServeMux) {
	s := r.Service
	mux.Handle("GET /v1/study/courses", MyCoursesHandler{service: s})
	mux.Handle("POST /v1/study/courses", CreateCourseHandler{service: s})
	mux.Handle("GET /v1/study/catalog", CatalogHandler{service: s})
	mux.Handle("GET /v1/study/courses/{id}", CourseHandler{service: s})
	mux.Handle("PUT /v1/study/courses/{id}", UpdateCourseHandler{service: s})
	mux.Handle("POST /v1/study/courses/{id}/archive", ArchiveCourseHandler{service: s, archived: true})
	mux.Handle("POST /v1/study/courses/{id}/restore", ArchiveCourseHandler{service: s})
	mux.Handle("GET /v1/study/courses/{id}/students", StudentsHandler{service: s})
	for _, list := range []List{ListTeachers, ListUnits, ListStudents} {
		path := "/v1/study/courses/{id}/" + string(list) + "/{member}"
		mux.Handle("PUT "+path, RosterChangeHandler{service: s, list: list, add: true})
		mux.Handle("DELETE "+path, RosterChangeHandler{service: s, list: list})
	}
}
